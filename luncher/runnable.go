package luncher

type Runnable interface {
	Run() error
	Shutdown() error
}
