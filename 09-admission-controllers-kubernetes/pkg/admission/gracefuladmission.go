package admission

import (
	"encoding/json"
	"fmt"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

func (g GracefulAdmissionHandler) AdmitValidation(r *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error) {
	p, err := requestToPod(r)
	if err != nil {
		fmt.Errorf("Could not get pod body from request")
		return nil, err
	}

	val, err := Validate(p)

	if err != nil {
		e := fmt.Sprintf("could not validate pod: %v", err)
		return generateAdmissionResponse(r.UID, false, 400, e), err
	}

	if !val.Valid {
		e := fmt.Sprintf("pod spec is invalid")
		return generateAdmissionResponse(r.UID, false, 400, e), nil
	}

	return generateAdmissionResponse(r.UID, true, 200, ""), nil

}

func (g GracefulAdmissionHandler) AdmitMutation(r *admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error) {
	p, err := requestToPod(r)
	if err != nil {
		fmt.Errorf("Could not get pod body from request")
		return nil, err
	}

	patches, err := MutatePatch(p)

	if err != nil {
		e := fmt.Sprintf("couldnt generate patch")
		return generateAdmissionResponse(r.UID, false, 400, e), nil
	}

	patcheb, err := json.Marshal(patches)

	if err != nil {
		fmt.Errorf("could not marshal JSON patch: %v", err)
	}

	return generatePatchAdmissionResponse(r.UID, patcheb)
}

func Validate(p *corev1.Pod) (ValidationResult, error) {
	if p.Spec.TerminationGracePeriodSeconds == nil {
		return ValidationResult{Valid: false, Reason: "no termination grace period"}, nil
	}
	for _, c := range p.Spec.Containers {
		if c.Lifecycle.PreStop == nil {
			return ValidationResult{Valid: false, Reason: "no pre stop hook"}, nil
		}
	}
	return ValidationResult{Valid: true, Reason: "valid"}, nil
}

func MutatePatch(p *corev1.Pod) ([]PatchOp, error) {
	defTgps := 60
	defPreStop := &corev1.LifecycleHandler{Sleep: &corev1.SleepAction{Seconds: 10}}
	patches := []PathOp

	if p.Spec.TerminationGracePeriodSeconds == nil {
		patches = append(patches, PatchOp{Op: "add", Path: "", Value: defTgps})
	}

	for _, c := range p.Spec.Containers {
		if c.Lifecycle.PreStop == nil {
			lc := c.Lifecycle.DeepCopy()
			lc.PreStop = defPreStop
			patches = append(patches, PatchOp{Op: "add", Path: "", Value: lc})
		}
	}

	return patches, nil

}
