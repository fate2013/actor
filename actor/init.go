package actor

var (
	fae *FaeExecutor
)

func init() {
	fae = NewFaeExecutor()
	fae.StartCluster()
}
