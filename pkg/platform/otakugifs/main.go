package otakugifs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{}

type reaction string

const (
	ReactionAirKiss    reaction = "airkiss"
	ReactionAngryStare reaction = "angrystare"
	ReactionBite       reaction = "bite"
	ReactionBleh       reaction = "bleh"
	ReactionBlush      reaction = "blush"
	ReactionBroFist    reaction = "brofist"
	ReactionCelebrate  reaction = "celebrate"
	ReactionCheers     reaction = "cheers"
	ReactionClap       reaction = "clap"
	ReactionConfused   reaction = "confused"
	ReactionCool       reaction = "cool"
	ReactionCry        reaction = "cry"
	ReactionCuddle     reaction = "cuddle"
	ReactionDance      reaction = "dance"
	ReactionDrool      reaction = "drool"
	ReactionEvilLaugh  reaction = "evillaugh"
	ReactionFacePalm   reaction = "facepalm"
	ReactionHandHold   reaction = "handhold"
	ReactionHappy      reaction = "happy"
	ReactionHeadBang   reaction = "headbang"
	ReactionHug        reaction = "hug"
	ReactionHuh        reaction = "huh"
	ReactionKiss       reaction = "kiss"
	ReactionLaugh      reaction = "laugh"
	ReactionLick       reaction = "lick"
	ReactionLove       reaction = "love"
	ReactionMad        reaction = "mad"
	ReactionNervous    reaction = "nervous"
	ReactionNo         reaction = "no"
	ReactionNom        reaction = "nom"
	ReactionNoseBleed  reaction = "nosebleed"
	ReactionNuzzle     reaction = "nuzzle"
	ReactionNyah       reaction = "nyah"
	ReactionPat        reaction = "pat"
	ReactionPeek       reaction = "peek"
	ReactionPinch      reaction = "pinch"
	ReactionPoke       reaction = "poke"
	ReactionPout       reaction = "pout"
	ReactionPunch      reaction = "punch"
	ReactionRoll       reaction = "roll"
	ReactionRun        reaction = "run"
	ReactionSad        reaction = "sad"
	ReactionScared     reaction = "scared"
	ReactionShout      reaction = "shout"
	ReactionShrug      reaction = "shrug"
	ReactionShy        reaction = "shy"
	ReactionSigh       reaction = "sigh"
	ReactionSing       reaction = "sing"
	ReactionSip        reaction = "sip"
	ReactionSlap       reaction = "slap"
	ReactionSleep      reaction = "sleep"
	ReactionSlowClap   reaction = "slowclap"
	ReactionSmack      reaction = "smack"
	ReactionSmile      reaction = "smile"
	ReactionSmug       reaction = "smug"
	ReactionSneeze     reaction = "sneeze"
	ReactionSorry      reaction = "sorry"
	ReactionStare      reaction = "stare"
	ReactionStop       reaction = "stop"
	ReactionSurprised  reaction = "surprised"
	ReactionSweat      reaction = "sweat"
	ReactionThumbsUp   reaction = "thumbsup"
	ReactionTickle     reaction = "tickle"
	ReactionTired      reaction = "tired"
	ReactionWave       reaction = "wave"
	ReactionWink       reaction = "wink"
	ReactionWoah       reaction = "woah"
	ReactionYawn       reaction = "yawn"
	ReactionYay        reaction = "yay"
	ReactionYes        reaction = "yes"
)

type format string

const (
	formatGIF  format = "GIF"
	formatWebP format = "WebP"
	formatAVIF format = "AVIF"
)

type resultStruct struct {
	URL string `json:"url"`
}

const imageURL = "https://api.otakugifs.xyz/gif?reaction=%s&format=%s"

func fetchImage(reaction reaction, format format) (ImageURL string, err error) {
	res, err := client.Get(fmt.Sprintf(imageURL, reaction, format))

	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	result := resultStruct{}

	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if result.URL == "" {
		return "", errors.New("An error occured fetching the image.")
	}

	return result.URL, nil
}

func FetchGIF(reaction reaction) (ImageURL string, err error) {
	return fetchImage(reaction, formatGIF)
}

func FetchWebP(reaction reaction) (ImageURL string, err error) {
	return fetchImage(reaction, formatWebP)
}

func FetchAVIF(reaction reaction) (ImageURL string, err error) {
	return fetchImage(reaction, formatAVIF)
}
