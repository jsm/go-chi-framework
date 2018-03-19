package validate

import "strings"

var profanity = []string{
	"abbo",
	"abo",
	"abuse",
	"adult",
	"alla",
	"allah",
	"anal",
	"angie",
	"angry",
	"anus",
	"arab",
	"arabs",
	"argie",
	"arse",
	"asian",
	"ass",
	"asses",
	"babe",
	"balls",
	"barf",
	"bast",
	"beast",
	"bi",
	"bible",
	"bitch",
	"black",
	"blind",
	"blow",
	"boang",
	"bogan",
	"bomb",
	"bombs",
	"bomd",
	"boner",
	"bong",
	"boob",
	"boobs",
	"booby",
	"boody",
	"boom",
	"boong",
	"booty",
	"bra",
	"bunga",
	"burn",
	"butt",
	"chav",
	"chin",
	"chink",
	"choad",
	"chode",
	"cigs",
	"clit",
	"cock",
	"cocky",
	"cohee",
	"color",
	"cooly",
	"coon",
	"cra5h",
	"crabs",
	"crack",
	"crap",
	"crash",
	"crime",
	"cum",
	"cumm",
	"cunn",
	"cunt",
	"dago",
	"damn",
	"darky",
	"dead",
	"death",
	"dego",
	"demon",
	"deth",
	"devil",
	"dick",
	"die",
	"died",
	"dies",
	"dike",
	"dildo",
	"dink",
	"dirty",
	"dive",
	"dix",
	"dong",
	"doom",
	"dope",
	"drug",
	"drunk",
	"dumb",
	"dyke",
	"eatme",
	"enema",
	"enemy",
	"erect",
	"ero",
	"evl",
	"fag",
	"fagot",
	"fairy",
	"faith",
	"fart",
	"fat",
	"fatah",
	"fatso",
	"fear",
	"feces",
	"felch",
	"fight",
	"fire",
	"floo",
	"fok",
	"fore",
	"forni",
	"fraud",
	"fu",
	"fubar",
	"fuc",
	"fucck",
	"fuck",
	"fucka",
	"fucks",
	"fugly",
	"fuk",
	"fuks",
	"fuuck",
	"gay",
	"geez",
	"geni",
	"gin",
	"ginzo",
	"gipp",
	"girls",
	"gob",
	"god",
	"gook",
	"goy",
	"goyim",
	"groe",
	"gross",
	"gubba",
	"gun",
	"gyp",
	"gypo",
	"gypp",
	"gyppo",
	"gyppy",
	"hamas",
	"hapa",
	"harem",
	"hebe",
	"heeb",
	"hell",
	"hiv",
	"ho",
	"hobo",
	"hoes",
	"hole",
	"homo",
	"honk",
	"honky",
	"hook",
	"hore",
	"hork",
	"horn",
	"horny",
	"hoser",
	"husky",
	"hussy",
	"hymen",
	"hymie",
	"idiot",
	"itch",
	"jade",
	"jap",
	"jebus",
	"jeez",
	"jesus",
	"jew",
	"jiga",
	"jigg",
	"jigga",
	"jiggy",
	"jihad",
	"jism",
	"jiz",
	"jizim",
	"jizm",
	"jizz",
	"joint",
	"jugs",
	"kafir",
	"kike",
	"kill",
	"kills",
	"kink",
	"kinky",
	"kkk",
	"knife",
	"kock",
	"koon",
	"kotex",
	"krap",
	"kraut",
	"kum",
	"kums",
	"kunt",
	"ky",
	"kyke",
	"laid",
	"latin",
	"lesbo",
	"lez",
	"lezbe",
	"lezbo",
	"lezz",
	"lezzo",
	"lies",
	"limey",
	"limy",
	"loser",
	"lsd",
	"lugan",
	"lynch",
	"mafia",
	"mams",
	"meth",
	"mgger",
	"mggor",
	"milf",
	"mocky",
	"mofo",
	"moky",
	"moles",
	"moron",
	"muff",
	"munt",
	"naked",
	"nasty",
	"nazi",
	"necro",
	"negro",
	"nig",
	"niger",
	"nigg",
	"nigga",
	"nigr",
	"nigra",
	"nigre",
	"nip",
	"nook",
	"nude",
	"nuke",
	"nymph",
	"oral",
	"orga",
	"orgy",
	"osama",
	"paki",
	"pansy",
	"panti",
	"payo",
	"peck",
	"pee",
	"pendy",
	"peni5",
	"penis",
	"perv",
	"phuk",
	"phuq",
	"pi55",
	"piker",
	"pikey",
	"piky",
	"pimp",
	"piss",
	"pixie",
	"pixy",
	"pocha",
	"pocho",
	"pohm",
	"pom",
	"pommy",
	"poo",
	"poon",
	"poop",
	"porn",
	"porno",
	"pot",
	"pric",
	"prick",
	"pros",
	"pu55i",
	"pu55y",
	"pube",
	"pubic",
	"pud",
	"pudd",
	"puke",
	"puss",
	"pussy",
	"pusy",
	"queef",
	"queer",
	"quim",
	"ra8s",
	"rabbi",
	"randy",
	"rape",
	"raped",
	"raper",
	"rere",
	"roach",
	"rump",
	"sadis",
	"sadom",
	"sandm",
	"satan",
	"scag",
	"scat",
	"screw",
	"scum",
	"semen",
	"seppo",
	"sex",
	"sexed",
	"sexy",
	"shag",
	"shat",
	"shav",
	"shhit",
	"shit",
	"shite",
	"shits",
	"shoot",
	"sick",
	"sissy",
	"skank",
	"skum",
	"slant",
	"slav",
	"slave",
	"slime",
	"slopy",
	"slut",
	"sluts",
	"slutt",
	"smack",
	"smut",
	"snot",
	"sob",
	"sodom",
	"sooty",
	"sos",
	"spank",
	"sperm",
	"spic",
	"spick",
	"spig",
	"spik",
	"spit",
	"spunk",
	"squaw",
	"stagg",
	"suck",
	"taboo",
	"taff",
	"tang",
	"tard",
	"teat",
	"teste",
	"tit",
	"tits",
	"titty",
	"tnt",
	"tramp",
	"trots",
	"turd",
	"twat",
	"twink",
	"uck",
	"uk",
	"urine",
	"usama",
	"vibr",
	"vomit",
	"vulva",
	"wab",
	"wank",
	"wetb",
	"whash",
	"whit",
	"whiz",
	"whop",
	"whore",
	"willy",
	"wn",
	"wog",
	"wop",
	"wtf",
	"wuss",
	"xtc",
	"xxx",
}

func IsProfane(inputString string) bool {
	for _, profaneWord := range profanity {
		if strings.Contains(inputString, profaneWord) {
			return true
		}
	}

	return false
}
