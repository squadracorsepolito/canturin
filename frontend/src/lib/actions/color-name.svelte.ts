import { darken, getLuminance, lighten } from 'color2k';
import randomColor from 'randomcolor';

type Options = {
	name: string;
};

export function colorByName(el: HTMLElement, opts: Options) {
	$effect(() => {
		const bgColor = randomColor({
			seed: opts.name
		});

		let textColor = '';
		if (getLuminance(bgColor) > 0.4) {
			textColor = darken(bgColor, 0.7);
		} else {
			textColor = lighten(bgColor, 0.7);
		}

		const borderColor = darken(bgColor, 0.15);

		el.style.backgroundColor = bgColor;
		el.style.color = textColor;
		el.style.borderColor = borderColor;
	});
}
