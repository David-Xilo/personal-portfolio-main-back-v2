package service

import (
	"math/rand"
	"personal-portfolio-main-back/src/internal/models"
	"time"
)

const ExtraReviews = 5

type PersonalReviewService struct {
	rng *rand.Rand
}

func NewPersonalReviewService() *PersonalReviewService {
	src := rand.NewSource(time.Now().UnixNano())
	return &PersonalReviewService{
		rng: rand.New(src),
	}
}

func (rs *PersonalReviewService) GetAllReviews() []*models.PersonalReviewsCarouselDTO {
	numberExtraReviews := rs.rng.Intn(ExtraReviews)
	var allReviews []*models.PersonalReviewsCarouselDTO

	goodReviews := rs.getFiveStarReviews()
	randomReviews := rs.getRandomReviews(numberExtraReviews)
	allReviews = append(allReviews, goodReviews...)
	allReviews = append(allReviews, randomReviews...)

	return allReviews
}

func (rs *PersonalReviewService) getRandomReviews(numberReviews int) []*models.PersonalReviewsCarouselDTO {
	var randomReviews []*models.PersonalReviewsCarouselDTO

	for i := 0; i < numberReviews; i++ {
		rating := rs.rng.Intn(5) + 1
		authorAndDescList := models.ReviewsByRating[rating]
		numberAuthorsAndDesc := len(authorAndDescList)
		chosenIdx := rs.rng.Intn(numberAuthorsAndDesc)
		review := models.CreatePersonalReviewsCarouselDTO(rating, authorAndDescList[chosenIdx].Description, authorAndDescList[chosenIdx].Author)
		randomReviews = append(randomReviews, review)
	}

	return randomReviews
}

// Can't have my mom or girlfriend giving me less than 5 stars lol
func (rs *PersonalReviewService) getFiveStarReviews() []*models.PersonalReviewsCarouselDTO {
	var fiveStarReviews []*models.PersonalReviewsCarouselDTO

	fiveStarReviews = append(fiveStarReviews, models.CreatePersonalReviewsCarouselDTO(5, "The most handsome!", "Mother"))
	fiveStarReviews = append(fiveStarReviews, models.CreatePersonalReviewsCarouselDTO(5, "Sweetie", "Wife"))
	fiveStarReviews = append(fiveStarReviews, models.CreatePersonalReviewsCarouselDTO(5, "A pretty cool dude", "Best Friend"))
	return fiveStarReviews
}
