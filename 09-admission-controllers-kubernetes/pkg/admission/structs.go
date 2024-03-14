package admission

import (
	admissionv1 "k8s.io/api/admission/v1"
)

type AdmissionHandler interface {
	AdmitMutation(*admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error)
	AdmitValidation(*admissionv1.AdmissionRequest) (*admissionv1.AdmissionReview, error)
}

type GracefulAdmissionHandler struct{}

type ValidationResult struct {
	Valid  bool
	Reason string
}

type PatchOp struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}
