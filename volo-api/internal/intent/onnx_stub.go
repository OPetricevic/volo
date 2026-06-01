package intent

// onnx_stub.go — Stub implementation when ONNX Runtime is not available.
// On Linux with onnxruntime installed, replace this with onnx_linux.go
// that uses github.com/yalue/onnxruntime_go.

// initOnnxSession is a no-op on platforms without ONNX Runtime.
func initOnnxSession(_ *ModelParser, _ string) bool {
	return false
}

// runOnnxInference is a no-op on platforms without ONNX Runtime.
func runOnnxInference(_ *ModelParser, _ []int64, _ []int64) []float32 {
	return nil
}
