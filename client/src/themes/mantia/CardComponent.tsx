import * as React from "react";
import { CardNumber } from "../CardNumber";
import { CardProps } from "../CardProps";
import { Season } from "../Season";
import { stylesheet } from "astroturf";

import ayame1 from "../mantia/ayame/1.svg";
import ayame2 from "../mantia/ayame/2.svg";
import ayame3 from "../mantia/ayame/3.svg";
import ayame4 from "../mantia/ayame/4.svg";

import botan1 from "../mantia/botan/1.svg";
import botan2 from "../mantia/botan/2.svg";
import botan3 from "../mantia/botan/3.svg";
import botan4 from "../mantia/botan/4.svg";

import fuji1 from "../mantia/fuji/1.svg";
import fuji2 from "../mantia/fuji/2.svg";
import fuji3 from "../mantia/fuji/3.svg";
import fuji4 from "../mantia/fuji/4.svg";

import hagi1 from "../mantia/hagi/1.svg";
import hagi2 from "../mantia/hagi/2.svg";
import hagi3 from "../mantia/hagi/3.svg";
import hagi4 from "../mantia/hagi/4.svg";

import kiku1 from "../mantia/kiku/1.svg";
import kiku2 from "../mantia/kiku/2.svg";
import kiku3 from "../mantia/kiku/3.svg";
import kiku4 from "../mantia/kiku/4.svg";

import kiri1 from "../mantia/kiri/1.svg";
import kiri2 from "../mantia/kiri/2.svg";
import kiri3 from "../mantia/kiri/3.svg";
import kiri4 from "../mantia/kiri/4.svg";

import matsu1 from "../mantia/matsu/1.svg";
import matsu2 from "../mantia/matsu/2.svg";
import matsu3 from "../mantia/matsu/3.svg";
import matsu4 from "../mantia/matsu/4.svg";

import momiji1 from "../mantia/momiji/1.svg";
import momiji2 from "../mantia/momiji/2.svg";
import momiji3 from "../mantia/momiji/3.svg";
import momiji4 from "../mantia/momiji/4.svg";

import sakura1 from "../mantia/sakura/1.svg";
import sakura2 from "../mantia/sakura/2.svg";
import sakura3 from "../mantia/sakura/3.svg";
import sakura4 from "../mantia/sakura/4.svg";

import susuki1 from "../mantia/susuki/1.svg";
import susuki2 from "../mantia/susuki/2.svg";
import susuki3 from "../mantia/susuki/3.svg";
import susuki4 from "../mantia/susuki/4.svg";

import ume1 from "../mantia/ume/1.svg";
import ume2 from "../mantia/ume/2.svg";
import ume3 from "../mantia/ume/3.svg";
import ume4 from "../mantia/ume/4.svg";

import yanagi1 from "../mantia/yanagi/1.svg";
import yanagi2 from "../mantia/yanagi/2.svg";
import yanagi3 from "../mantia/yanagi/3.svg";
import yanagi4 from "../mantia/yanagi/4.svg";

const SVG_MAP: Record<Season, Record<CardNumber, string>> = {
	[Season.Ayame]: {
		[CardNumber.First]: ayame1,
		[CardNumber.Second]: ayame2,
		[CardNumber.Third]: ayame3,
		[CardNumber.Fourth]: ayame4,
	},
	[Season.Botan]: {
		[CardNumber.First]: botan1,
		[CardNumber.Second]: botan2,
		[CardNumber.Third]: botan3,
		[CardNumber.Fourth]: botan4,
	},
	[Season.Fuji]: {
		[CardNumber.First]: fuji1,
		[CardNumber.Second]: fuji2,
		[CardNumber.Third]: fuji3,
		[CardNumber.Fourth]: fuji4,
	},
	[Season.Hagi]: {
		[CardNumber.First]: hagi1,
		[CardNumber.Second]: hagi2,
		[CardNumber.Third]: hagi3,
		[CardNumber.Fourth]: hagi4,
	},
	[Season.Kiku]: {
		[CardNumber.First]: kiku1,
		[CardNumber.Second]: kiku2,
		[CardNumber.Third]: kiku3,
		[CardNumber.Fourth]: kiku4,
	},
	[Season.Kiri]: {
		[CardNumber.First]: kiri1,
		[CardNumber.Second]: kiri2,
		[CardNumber.Third]: kiri3,
		[CardNumber.Fourth]: kiri4,
	},
	[Season.Matsu]: {
		[CardNumber.First]: matsu1,
		[CardNumber.Second]: matsu2,
		[CardNumber.Third]: matsu3,
		[CardNumber.Fourth]: matsu4,
	},
	[Season.Momiji]: {
		[CardNumber.First]: momiji1,
		[CardNumber.Second]: momiji2,
		[CardNumber.Third]: momiji3,
		[CardNumber.Fourth]: momiji4,
	},
	[Season.Sakura]: {
		[CardNumber.First]: sakura1,
		[CardNumber.Second]: sakura2,
		[CardNumber.Third]: sakura3,
		[CardNumber.Fourth]: sakura4,
	},
	[Season.Susuki]: {
		[CardNumber.First]: susuki1,
		[CardNumber.Second]: susuki2,
		[CardNumber.Third]: susuki3,
		[CardNumber.Fourth]: susuki4,
	},
	[Season.Ume]: {
		[CardNumber.First]: ume1,
		[CardNumber.Second]: ume2,
		[CardNumber.Third]: ume3,
		[CardNumber.Fourth]: ume4,
	},
	[Season.Yanagi]: {
		[CardNumber.First]: yanagi1,
		[CardNumber.Second]: yanagi2,
		[CardNumber.Third]: yanagi3,
		[CardNumber.Fourth]: yanagi4,
	}
}

const styles = stylesheet`
	.card {
		width: 100px;
	}
`;

export const CardComponent: React.FC<CardProps> = (props) => {
	return <div className={styles.card}>
		<img src={SVG_MAP[props.season][props.card]} />
	</div>
}
