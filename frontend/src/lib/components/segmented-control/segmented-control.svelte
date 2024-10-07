<script lang="ts">
	import * as radioGroup from '@zag-js/radio-group';
	import { useMachine, normalizeProps } from '@zag-js/svelte';
	import { TickIcon } from '../icon';
	import { uniqueId } from '$lib/utils';
	import type { SegmentedControlOption } from './types';

	type Props = {
		selectedValue: string;
		name: string;
		options: SegmentedControlOption[];
		readOnly?: boolean;
	};

	let { selectedValue = $bindable(), name, options, readOnly }: Props = $props();

	const [snapshot, send] = useMachine(
		radioGroup.machine({
			id: uniqueId(),
			name,
			orientation: 'horizontal',
			value: selectedValue,
			readOnly,
			onValueChange: (details) => {
				selectedValue = details.value;
			}
		})
	);

	const api = $derived(radioGroup.connect(snapshot, send, normalizeProps));
</script>

<div
	{...api.getRootProps()}
	class="root"
	style:grid-template-columns="repeat({options.length}, minmax(0, 1fr))"
>
	{#each options as opt}
		<label {...api.getItemProps({ value: opt.value })} class="item">
			<div {...api.getItemTextProps({ value: opt.value })} class="label">
				<span>
					{opt.label}
				</span>

				{#if opt.value === api.value}
					<span>
						<TickIcon height={20} width={20} />
					</span>
				{/if}
			</div>

			{#if opt.desc}
				<span class="desc">{opt.desc}</span>
			{/if}

			<input {...api.getItemHiddenInputProps({ value: opt.value })} />
		</label>
	{/each}
</div>

<style lang="postcss">
	.root {
		@apply grid gap-x-2;

		.item {
			@apply grid row-span-2 grid-rows-subgrid border-2 py-1 px-2 rounded-btn gap-1 transition-colors;

			&[data-readonly] {
				&:not([data-state='checked']) {
					@apply bg-base-300 border-base-300 opacity-80;
				}
			}

			&:not([data-readonly]) {
				@apply border-base-300 cursor-pointer;

				&:not([data-state='checked']) {
					&:hover {
						@apply border-secondary text-secondary bg-secondary-ghost;
					}
				}

				&[data-focus] {
					@apply focus-ring-primary;
				}
			}

			&[data-state='checked'] {
				@apply border-primary text-primary bg-primary-ghost;
			}

			.label {
				@apply font-medium flex items-center gap-1 justify-between;
			}

			.desc {
				@apply text-sm text-dimmed;
			}
		}
	}
</style>
