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

	const [snapshot, send] = useMachine(
		zagSwitch.machine({
			id: uniqueId(),

			onCheckedChange: (details) => {
				checked = details.checked;
				oncheckchange?.(details.checked);
			}
		}),
		{
			context: {
				get checked() {
					return checked;
				},
				get readOnly() {
					return readOnly;
				}
			}
		}
	);

	const api = $derived(zagSwitch.connect(snapshot, send, normalizeProps));
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
	[data-part='control'] {
		@apply flex border-2 w-12 p-1 rounded-badge transition-colors;

		&:not([data-readonly]) {
			@apply cursor-pointer;
		}

		&[data-state='checked'] {
			@apply border-primary justify-end;
		}

		&[data-state='unchecked'] {
			@apply border-base-300 justify-start;

			&[data-readonly] {
				@apply bg-base-300 opacity-80;
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
			@apply bg-base-content;
		}
	}
</style>
