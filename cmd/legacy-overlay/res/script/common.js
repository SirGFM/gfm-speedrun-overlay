// /script/common.js contain some ugly constants shared by both control and view.

const DIRATTR = 'data-direction';
const ATTR = 'data-updated-field';
const WINNERS = "winners";
const LOSERS = "losers";
const POOLS = "pools";
const SEMIS = "semis";
const FINALS = "finals";
const GRANDFINALS = "grand-finals";
const BO3 = 3;
const BO5 = 5;
const RESET = "Reset";
const DEC = "-";
const INC = "+";

var state = {
    game: {
        tournament: "NAME",
        side: WINNERS,
        numMatches: BO3,
        round: POOLS,
        rpool: 1,
        match: 1,
    },
    left: {
        player: -1,
        character: -1,
        score: 0
    },
    right: {
        player: -1,
        character: -1,
        score: 0
    },
    other: {
        camera: true,
        gamepad: true
    }
};

var character_list = [
    {
        "name": "Absa",
        "img": "/img/char/absa.png",
        "tag": "/img/tag/absa.png"
    },
    {
        "name": "Clairen",
        "img": "/img/char/clairen.png",
        "tag": "/img/tag/clairen.png"
    },
    {
        "name": "Etalus",
        "img": "/img/char/etalus.png",
        "tag": "/img/tag/etalus.png"
    },
    {
        "name": "Forsburn",
        "img": "/img/char/forsburn.png",
        "tag": "/img/tag/forsburn.png"
    },
    {
        "name": "Kragg",
        "img": "/img/char/kragg.png",
        "tag": "/img/tag/kragg.png"
    },
    {
        "name": "Maypul",
        "img": "/img/char/maypul.png",
        "tag": "/img/tag/maypul.png"
    },
    {
        "name": "Orcane",
        "img": "/img/char/orcane.png",
        "tag": "/img/tag/orcane.png"
    },
    {
        "name": "Ori",
        "img": "/img/char/ori.png",
        "tag": "/img/tag/ori.png"
    },
    {
        "name": "Ranno",
        "img": "/img/char/ranno.png",
        "tag": "/img/tag/ranno.png"
    },
    {
        "name": "Wrastor",
        "img": "/img/char/wrastor.png",
        "tag": "/img/tag/wrastor.png"
    },
    {
        "name": "Zetterburn",
        "img": "/img/char/zetterburn.png",
        "tag": "/img/tag/zetterburn.png"
    }
];
