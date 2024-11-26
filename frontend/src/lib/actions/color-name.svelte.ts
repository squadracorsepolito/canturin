import { getColorByName } from '$lib/utils';

type Options = {
	name: string;
};

export function colorByName(el: HTMLElement, opts: Options) {
	$effect(() => {
		const colors = getColorByName(opts.name);

		el.style.backgroundColor = colors.bgColor;
		el.style.color = colors.textColor;
	});
}
