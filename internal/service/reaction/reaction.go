package reaction

type IReactionService interface {
}

type reactionService struct {
	reactsRepo reaction.IReactionsRepository
}

func NewReactionService()
