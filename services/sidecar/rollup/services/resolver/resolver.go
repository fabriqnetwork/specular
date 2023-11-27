package resolver

type ResolverService struct {
}

func (r *ResolverService) Run() {

}

func (r *ResolverService) ConfirmAssertion() {
	// If the first unresolved assertion is eligible for confirmation, trigger its confirmation. Otherwise, wait.

}

func (r *ResolverService) RejectAssertion() {

}
