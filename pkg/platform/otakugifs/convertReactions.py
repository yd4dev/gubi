import json
import urllib.request

url = "https://api.otakugifs.xyz/gif/allreactions"

# Fetch the JSON from the API
with urllib.request.urlopen(url) as response:
    data = response.read()
    reaction = json.loads(data)

camel_case_overrides = {
    "airkiss": "AirKiss",
    "angrystare": "AngryStare",
    "brofist": "BroFist",
    "evillaugh": "EvilLaugh",
    "facepalm": "FacePalm",
    "handhold": "HandHold",
    "headbang": "HeadBang",
    "nosebleed": "NoseBleed",
    "slowclap": "SlowClap",
    "thumbsup": "ThumbsUp",
}

for r in reaction["reactions"]:
    # Capitalize using the overrides dictionary, or fallback to simple capitalization
    name = camel_case_overrides.get(r, r.capitalize())

    print(f'Reaction{name} Reaction = "{r}"')
