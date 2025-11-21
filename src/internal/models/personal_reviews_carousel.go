package models

type PersonalReviewsCarouselDTO struct {
	Author      string `json:"author"`
	Description string `json:"description"`
	Rating      int    `json:"rating"`
}

func CreatePersonalReviewsCarouselDTO(rating int, desc string, author string) *PersonalReviewsCarouselDTO {
	return &PersonalReviewsCarouselDTO{
		Author:      author,
		Description: desc,
		Rating:      rating,
	}
}

type ReviewTuple struct {
	Author      string
	Description string
}

var ReviewsByRating = map[int][]ReviewTuple{
	5: {
		{"Dog", "Oof oof!"},
		{"Grandparents", "Such a good boy"},
	},
	4: {
		{"Siblings", "He can be pretty tchill"},
		{"Swimming Coach", "Not a great swimmer but he's a nice person"},
	},
	3: {
		{"Neighbor", "He can be a bit loud, but nice overall"},
		{"Third Grade Teacher", "He was kind of a pain"},
	},
	2: {
		{"Childhood Rival", "I'll beat him someday"},
		{"Time Traveler", "Keeps preventing the disasters I'm trying to visit"},
	},
	1: {
		{"James Bond Villain", "I'll get you next time!"},
		{"Alien", "Keeps dodging our beam"},
	},
}
