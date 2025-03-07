<script lang="ts">
	import { normalizeProps, useMachine, mergeProps } from '@zag-js/svelte';
	import * as checkbox from '@zag-js/checkbox';
	import { uniqueId } from '$lib/utils';

	type Props = {
		checked: boolean;
		label?: string;
		oncheckchange?: (checked: boolean) => void;
	};

	let { checked = $bindable(), label, oncheckchange }: Props = $props();

	const checkboxProps: checkbox.Props = $derived({
		id: uniqueId(),
		checked: checked,
		onCheckedChange: (details) => {
			let tmpChecked = false;
			if (details.checked !== 'indeterminate') {
				tmpChecked = details.checked;
			}

			checked = tmpChecked;
			oncheckchange?.(tmpChecked);
		}
	});

	const service = useMachine(checkbox.machine, () => checkboxProps);

	const api = $derived(checkbox.connect(service, normalizeProps));

	const controlProps = $derived(
		mergeProps(api.getControlProps(), {
			'aria-checked': api.checked
		})
	);
</script>

<label {...api.getRootProps()}>
	{#if label}
		<span {...api.getLabelProps()}>
			{label}
		</span>
	{/if}
	<div {...controlProps}></div>
	<input {...api.getHiddenInputProps()} />
</label>

<style lang="postcss">
	[data-part='root'] {
		@apply inline-block;
	}

	[data-part='control'] {
		@apply checkbox h-5 w-5;
	}
</style>
