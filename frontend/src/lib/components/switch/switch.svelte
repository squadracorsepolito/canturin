<script lang="ts">
	import { uniqueId } from '$lib/utils';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import * as zagSwitch from '@zag-js/switch';

	type Props = {
		checked: boolean;
		label?: string;
		readOnly?: boolean;
		oncheckedchange?: (checked: boolean) => void;
	};

	let { checked = $bindable(), label, readOnly, oncheckedchange: oncheckchange }: Props = $props();

	const switchProps: zagSwitch.Props = $derived({
		id: uniqueId(),
		checked: checked,
		readOnly: readOnly,
		onCheckedChange: (details) => {
			checked = details.checked;
			oncheckchange?.(details.checked);
		}
	});

	const service = useMachine(zagSwitch.machine, () => switchProps);

	const api = $derived(zagSwitch.connect(service, normalizeProps));
</script>

<label {...api.getRootProps()}>
	<input {...api.getHiddenInputProps()} />
	<span {...api.getControlProps()}>
		<span {...api.getThumbProps()}></span>
	</span>
	{#if label}
		<span {...api.getLabelProps()}>{label}</span>
	{/if}
</label>

<style lang="postcss">
	[data-part='root'] {
		@apply inline-block;
	}

	[data-part='control'] {
		@apply flex border-2 w-12 p-1 rounded-badge transition-colors;

		&:not([data-readonly]) {
			@apply cursor-pointer;
		}

		&[data-state='checked'] {
			@apply border-primary justify-end;
		}

		&[data-state='unchecked'] {
			@apply border-neutral-content justify-start;

			&[data-readonly] {
				@apply bg-neutral-content opacity-80;
			}
		}

		&[data-focus] {
			@apply focus-ring-primary;
		}
	}

	[data-part='thumb'] {
		@apply block rounded-full h-4 w-4;

		&[data-state='checked'] {
			@apply bg-primary;
		}

		&[data-state='unchecked'] {
			@apply bg-neutral;
		}
	}
</style>
