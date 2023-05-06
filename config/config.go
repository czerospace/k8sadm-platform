package config

const (
	ListenAddr  = "0.0.0.0:9090"
	Kubeconfigs = `{"TST-1":"config/opsconfig","TST-2":"config/opsconfig"}`

	// PodLogTailLine 查看容器日志时，显示的 tail 行数 tail -n 5000
	PodLogTailLine = 5000
)
