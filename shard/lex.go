package shard

// These are the numbers 0 to 1000 encoded as strings that sort lexicographically.
// We use these to sort on the datastore key of the sharded value (we can't just use
// the raw integer values because they don't sort the right way).
var positionKeys = []string{
	"A0",
	"A1",
	"A2",
	"A3",
	"A4",
	"A5",
	"A6",
	"A7",
	"A8",
	"A9",
	"AA",
	"AB",
	"AC",
	"AD",
	"AE",
	"AF",
	"AG",
	"AH",
	"AI",
	"AJ",
	"AK",
	"AL",
	"AM",
	"AN",
	"AO",
	"AP",
	"AQ",
	"AR",
	"AS",
	"AT",
	"AU",
	"AV",
	"AW",
	"AX",
	"AY",
	"AZ",
	"Aa",
	"Ab",
	"Ac",
	"Ad",
	"Ae",
	"Af",
	"Ag",
	"Ah",
	"Ai",
	"Aj",
	"Ak",
	"Al",
	"Am",
	"An",
	"Ao",
	"Ap",
	"Aq",
	"Ar",
	"As",
	"At",
	"Au",
	"Av",
	"Aw",
	"Ax",
	"Ay",
	"Az",
	"B10",
	"B11",
	"B12",
	"B13",
	"B14",
	"B15",
	"B16",
	"B17",
	"B18",
	"B19",
	"B1A",
	"B1B",
	"B1C",
	"B1D",
	"B1E",
	"B1F",
	"B1G",
	"B1H",
	"B1I",
	"B1J",
	"B1K",
	"B1L",
	"B1M",
	"B1N",
	"B1O",
	"B1P",
	"B1Q",
	"B1R",
	"B1S",
	"B1T",
	"B1U",
	"B1V",
	"B1W",
	"B1X",
	"B1Y",
	"B1Z",
	"B1a",
	"B1b",
	"B1c",
	"B1d",
	"B1e",
	"B1f",
	"B1g",
	"B1h",
	"B1i",
	"B1j",
	"B1k",
	"B1l",
	"B1m",
	"B1n",
	"B1o",
	"B1p",
	"B1q",
	"B1r",
	"B1s",
	"B1t",
	"B1u",
	"B1v",
	"B1w",
	"B1x",
	"B1y",
	"B1z",
	"B20",
	"B21",
	"B22",
	"B23",
	"B24",
	"B25",
	"B26",
	"B27",
	"B28",
	"B29",
	"B2A",
	"B2B",
	"B2C",
	"B2D",
	"B2E",
	"B2F",
	"B2G",
	"B2H",
	"B2I",
	"B2J",
	"B2K",
	"B2L",
	"B2M",
	"B2N",
	"B2O",
	"B2P",
	"B2Q",
	"B2R",
	"B2S",
	"B2T",
	"B2U",
	"B2V",
	"B2W",
	"B2X",
	"B2Y",
	"B2Z",
	"B2a",
	"B2b",
	"B2c",
	"B2d",
	"B2e",
	"B2f",
	"B2g",
	"B2h",
	"B2i",
	"B2j",
	"B2k",
	"B2l",
	"B2m",
	"B2n",
	"B2o",
	"B2p",
	"B2q",
	"B2r",
	"B2s",
	"B2t",
	"B2u",
	"B2v",
	"B2w",
	"B2x",
	"B2y",
	"B2z",
	"B30",
	"B31",
	"B32",
	"B33",
	"B34",
	"B35",
	"B36",
	"B37",
	"B38",
	"B39",
	"B3A",
	"B3B",
	"B3C",
	"B3D",
	"B3E",
	"B3F",
	"B3G",
	"B3H",
	"B3I",
	"B3J",
	"B3K",
	"B3L",
	"B3M",
	"B3N",
	"B3O",
	"B3P",
	"B3Q",
	"B3R",
	"B3S",
	"B3T",
	"B3U",
	"B3V",
	"B3W",
	"B3X",
	"B3Y",
	"B3Z",
	"B3a",
	"B3b",
	"B3c",
	"B3d",
	"B3e",
	"B3f",
	"B3g",
	"B3h",
	"B3i",
	"B3j",
	"B3k",
	"B3l",
	"B3m",
	"B3n",
	"B3o",
	"B3p",
	"B3q",
	"B3r",
	"B3s",
	"B3t",
	"B3u",
	"B3v",
	"B3w",
	"B3x",
	"B3y",
	"B3z",
	"B40",
	"B41",
	"B42",
	"B43",
	"B44",
	"B45",
	"B46",
	"B47",
	"B48",
	"B49",
	"B4A",
	"B4B",
	"B4C",
	"B4D",
	"B4E",
	"B4F",
	"B4G",
	"B4H",
	"B4I",
	"B4J",
	"B4K",
	"B4L",
	"B4M",
	"B4N",
	"B4O",
	"B4P",
	"B4Q",
	"B4R",
	"B4S",
	"B4T",
	"B4U",
	"B4V",
	"B4W",
	"B4X",
	"B4Y",
	"B4Z",
	"B4a",
	"B4b",
	"B4c",
	"B4d",
	"B4e",
	"B4f",
	"B4g",
	"B4h",
	"B4i",
	"B4j",
	"B4k",
	"B4l",
	"B4m",
	"B4n",
	"B4o",
	"B4p",
	"B4q",
	"B4r",
	"B4s",
	"B4t",
	"B4u",
	"B4v",
	"B4w",
	"B4x",
	"B4y",
	"B4z",
	"B50",
	"B51",
	"B52",
	"B53",
	"B54",
	"B55",
	"B56",
	"B57",
	"B58",
	"B59",
	"B5A",
	"B5B",
	"B5C",
	"B5D",
	"B5E",
	"B5F",
	"B5G",
	"B5H",
	"B5I",
	"B5J",
	"B5K",
	"B5L",
	"B5M",
	"B5N",
	"B5O",
	"B5P",
	"B5Q",
	"B5R",
	"B5S",
	"B5T",
	"B5U",
	"B5V",
	"B5W",
	"B5X",
	"B5Y",
	"B5Z",
	"B5a",
	"B5b",
	"B5c",
	"B5d",
	"B5e",
	"B5f",
	"B5g",
	"B5h",
	"B5i",
	"B5j",
	"B5k",
	"B5l",
	"B5m",
	"B5n",
	"B5o",
	"B5p",
	"B5q",
	"B5r",
	"B5s",
	"B5t",
	"B5u",
	"B5v",
	"B5w",
	"B5x",
	"B5y",
	"B5z",
	"B60",
	"B61",
	"B62",
	"B63",
	"B64",
	"B65",
	"B66",
	"B67",
	"B68",
	"B69",
	"B6A",
	"B6B",
	"B6C",
	"B6D",
	"B6E",
	"B6F",
	"B6G",
	"B6H",
	"B6I",
	"B6J",
	"B6K",
	"B6L",
	"B6M",
	"B6N",
	"B6O",
	"B6P",
	"B6Q",
	"B6R",
	"B6S",
	"B6T",
	"B6U",
	"B6V",
	"B6W",
	"B6X",
	"B6Y",
	"B6Z",
	"B6a",
	"B6b",
	"B6c",
	"B6d",
	"B6e",
	"B6f",
	"B6g",
	"B6h",
	"B6i",
	"B6j",
	"B6k",
	"B6l",
	"B6m",
	"B6n",
	"B6o",
	"B6p",
	"B6q",
	"B6r",
	"B6s",
	"B6t",
	"B6u",
	"B6v",
	"B6w",
	"B6x",
	"B6y",
	"B6z",
	"B70",
	"B71",
	"B72",
	"B73",
	"B74",
	"B75",
	"B76",
	"B77",
	"B78",
	"B79",
	"B7A",
	"B7B",
	"B7C",
	"B7D",
	"B7E",
	"B7F",
	"B7G",
	"B7H",
	"B7I",
	"B7J",
	"B7K",
	"B7L",
	"B7M",
	"B7N",
	"B7O",
	"B7P",
	"B7Q",
	"B7R",
	"B7S",
	"B7T",
	"B7U",
	"B7V",
	"B7W",
	"B7X",
	"B7Y",
	"B7Z",
	"B7a",
	"B7b",
	"B7c",
	"B7d",
	"B7e",
	"B7f",
	"B7g",
	"B7h",
	"B7i",
	"B7j",
	"B7k",
	"B7l",
	"B7m",
	"B7n",
	"B7o",
	"B7p",
	"B7q",
	"B7r",
	"B7s",
	"B7t",
	"B7u",
	"B7v",
	"B7w",
	"B7x",
	"B7y",
	"B7z",
	"B80",
	"B81",
	"B82",
	"B83",
	"B84",
	"B85",
	"B86",
	"B87",
	"B88",
	"B89",
	"B8A",
	"B8B",
	"B8C",
	"B8D",
	"B8E",
	"B8F",
	"B8G",
	"B8H",
	"B8I",
	"B8J",
	"B8K",
	"B8L",
	"B8M",
	"B8N",
	"B8O",
	"B8P",
	"B8Q",
	"B8R",
	"B8S",
	"B8T",
	"B8U",
	"B8V",
	"B8W",
	"B8X",
	"B8Y",
	"B8Z",
	"B8a",
	"B8b",
	"B8c",
	"B8d",
	"B8e",
	"B8f",
	"B8g",
	"B8h",
	"B8i",
	"B8j",
	"B8k",
	"B8l",
	"B8m",
	"B8n",
	"B8o",
	"B8p",
	"B8q",
	"B8r",
	"B8s",
	"B8t",
	"B8u",
	"B8v",
	"B8w",
	"B8x",
	"B8y",
	"B8z",
	"B90",
	"B91",
	"B92",
	"B93",
	"B94",
	"B95",
	"B96",
	"B97",
	"B98",
	"B99",
	"B9A",
	"B9B",
	"B9C",
	"B9D",
	"B9E",
	"B9F",
	"B9G",
	"B9H",
	"B9I",
	"B9J",
	"B9K",
	"B9L",
	"B9M",
	"B9N",
	"B9O",
	"B9P",
	"B9Q",
	"B9R",
	"B9S",
	"B9T",
	"B9U",
	"B9V",
	"B9W",
	"B9X",
	"B9Y",
	"B9Z",
	"B9a",
	"B9b",
	"B9c",
	"B9d",
	"B9e",
	"B9f",
	"B9g",
	"B9h",
	"B9i",
	"B9j",
	"B9k",
	"B9l",
	"B9m",
	"B9n",
	"B9o",
	"B9p",
	"B9q",
	"B9r",
	"B9s",
	"B9t",
	"B9u",
	"B9v",
	"B9w",
	"B9x",
	"B9y",
	"B9z",
	"BA0",
	"BA1",
	"BA2",
	"BA3",
	"BA4",
	"BA5",
	"BA6",
	"BA7",
	"BA8",
	"BA9",
	"BAA",
	"BAB",
	"BAC",
	"BAD",
	"BAE",
	"BAF",
	"BAG",
	"BAH",
	"BAI",
	"BAJ",
	"BAK",
	"BAL",
	"BAM",
	"BAN",
	"BAO",
	"BAP",
	"BAQ",
	"BAR",
	"BAS",
	"BAT",
	"BAU",
	"BAV",
	"BAW",
	"BAX",
	"BAY",
	"BAZ",
	"BAa",
	"BAb",
	"BAc",
	"BAd",
	"BAe",
	"BAf",
	"BAg",
	"BAh",
	"BAi",
	"BAj",
	"BAk",
	"BAl",
	"BAm",
	"BAn",
	"BAo",
	"BAp",
	"BAq",
	"BAr",
	"BAs",
	"BAt",
	"BAu",
	"BAv",
	"BAw",
	"BAx",
	"BAy",
	"BAz",
	"BB0",
	"BB1",
	"BB2",
	"BB3",
	"BB4",
	"BB5",
	"BB6",
	"BB7",
	"BB8",
	"BB9",
	"BBA",
	"BBB",
	"BBC",
	"BBD",
	"BBE",
	"BBF",
	"BBG",
	"BBH",
	"BBI",
	"BBJ",
	"BBK",
	"BBL",
	"BBM",
	"BBN",
	"BBO",
	"BBP",
	"BBQ",
	"BBR",
	"BBS",
	"BBT",
	"BBU",
	"BBV",
	"BBW",
	"BBX",
	"BBY",
	"BBZ",
	"BBa",
	"BBb",
	"BBc",
	"BBd",
	"BBe",
	"BBf",
	"BBg",
	"BBh",
	"BBi",
	"BBj",
	"BBk",
	"BBl",
	"BBm",
	"BBn",
	"BBo",
	"BBp",
	"BBq",
	"BBr",
	"BBs",
	"BBt",
	"BBu",
	"BBv",
	"BBw",
	"BBx",
	"BBy",
	"BBz",
	"BC0",
	"BC1",
	"BC2",
	"BC3",
	"BC4",
	"BC5",
	"BC6",
	"BC7",
	"BC8",
	"BC9",
	"BCA",
	"BCB",
	"BCC",
	"BCD",
	"BCE",
	"BCF",
	"BCG",
	"BCH",
	"BCI",
	"BCJ",
	"BCK",
	"BCL",
	"BCM",
	"BCN",
	"BCO",
	"BCP",
	"BCQ",
	"BCR",
	"BCS",
	"BCT",
	"BCU",
	"BCV",
	"BCW",
	"BCX",
	"BCY",
	"BCZ",
	"BCa",
	"BCb",
	"BCc",
	"BCd",
	"BCe",
	"BCf",
	"BCg",
	"BCh",
	"BCi",
	"BCj",
	"BCk",
	"BCl",
	"BCm",
	"BCn",
	"BCo",
	"BCp",
	"BCq",
	"BCr",
	"BCs",
	"BCt",
	"BCu",
	"BCv",
	"BCw",
	"BCx",
	"BCy",
	"BCz",
	"BD0",
	"BD1",
	"BD2",
	"BD3",
	"BD4",
	"BD5",
	"BD6",
	"BD7",
	"BD8",
	"BD9",
	"BDA",
	"BDB",
	"BDC",
	"BDD",
	"BDE",
	"BDF",
	"BDG",
	"BDH",
	"BDI",
	"BDJ",
	"BDK",
	"BDL",
	"BDM",
	"BDN",
	"BDO",
	"BDP",
	"BDQ",
	"BDR",
	"BDS",
	"BDT",
	"BDU",
	"BDV",
	"BDW",
	"BDX",
	"BDY",
	"BDZ",
	"BDa",
	"BDb",
	"BDc",
	"BDd",
	"BDe",
	"BDf",
	"BDg",
	"BDh",
	"BDi",
	"BDj",
	"BDk",
	"BDl",
	"BDm",
	"BDn",
	"BDo",
	"BDp",
	"BDq",
	"BDr",
	"BDs",
	"BDt",
	"BDu",
	"BDv",
	"BDw",
	"BDx",
	"BDy",
	"BDz",
	"BE0",
	"BE1",
	"BE2",
	"BE3",
	"BE4",
	"BE5",
	"BE6",
	"BE7",
	"BE8",
	"BE9",
	"BEA",
	"BEB",
	"BEC",
	"BED",
	"BEE",
	"BEF",
	"BEG",
	"BEH",
	"BEI",
	"BEJ",
	"BEK",
	"BEL",
	"BEM",
	"BEN",
	"BEO",
	"BEP",
	"BEQ",
	"BER",
	"BES",
	"BET",
	"BEU",
	"BEV",
	"BEW",
	"BEX",
	"BEY",
	"BEZ",
	"BEa",
	"BEb",
	"BEc",
	"BEd",
	"BEe",
	"BEf",
	"BEg",
	"BEh",
	"BEi",
	"BEj",
	"BEk",
	"BEl",
	"BEm",
	"BEn",
	"BEo",
	"BEp",
	"BEq",
	"BEr",
	"BEs",
	"BEt",
	"BEu",
	"BEv",
	"BEw",
	"BEx",
	"BEy",
	"BEz",
	"BF0",
	"BF1",
	"BF2",
	"BF3",
	"BF4",
	"BF5",
	"BF6",
	"BF7",
	"BF8",
	"BF9",
	"BFA",
	"BFB",
	"BFC",
	"BFD",
	"BFE",
	"BFF",
	"BFG",
	"BFH",
	"BFI",
	"BFJ",
	"BFK",
	"BFL",
	"BFM",
	"BFN",
	"BFO",
	"BFP",
	"BFQ",
	"BFR",
	"BFS",
	"BFT",
	"BFU",
	"BFV",
	"BFW",
	"BFX",
	"BFY",
	"BFZ",
	"BFa",
	"BFb",
	"BFc",
	"BFd",
	"BFe",
	"BFf",
	"BFg",
	"BFh",
	"BFi",
	"BFj",
	"BFk",
	"BFl",
	"BFm",
	"BFn",
	"BFo",
	"BFp",
	"BFq",
	"BFr",
	"BFs",
	"BFt",
	"BFu",
	"BFv",
	"BFw",
	"BFx",
	"BFy",
	"BFz",
	"BG0",
	"BG1",
	"BG2",
	"BG3",
	"BG4",
	"BG5",
	"BG6",
	"BG7",
}