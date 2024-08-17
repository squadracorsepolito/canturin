import { darken, getLuminance, lighten } from 'color2k';
import randomColor from 'randomcolor';

export function getHexNumber(num: number) {
	return `0x${num.toString(16).padStart(2, '0')}`;
}

export function getColorByName(name: string) {
	const bgColor = randomColor({
		seed: name
	});

	let textColor = '';
	if (getLuminance(bgColor) > 0.4) {
		textColor = darken(bgColor, 0.7);
	} else {
		textColor = lighten(bgColor, 0.7);
	}

	return {
		bgColor: bgColor,
		textColor: textColor
	};
}
