<script lang="ts" generics="T extends { [K in keyof T]: any }">
	import { uniqueId, type KeyOfString } from '$lib/utils';
	import * as select from '@zag-js/select';
	import { normalizeProps, portal, useMachine } from '@zag-js/svelte';
	import { AltArrowIcon, CheckIcon } from '../icon';

	type Props = {
		items: T[];
		selected: T;
		name: string;
		valueKey: KeyOfString<T>;
		labelKey: KeyOfString<T>;
		onselect?: (item: T) => void;
		filter?: (item: T) => boolean;
	};

	let {
		items,
		selected = $bindable(),
		name,
		valueKey,
		labelKey,
		onselect,
		filter
	}: Props = $props();

	const collection = select.collection({
		items,
		isItemDisabled: filter,
		itemToString: (item) => item[labelKey],
		itemToValue: (item) => item[valueKey]
	});

	const [snapshot, send] = useMachine(
		select.machine({
			id: uniqueId(),
			collection,
			name: name,
			onValueChange: (details) => {
				if (details.items.length === 0) {
					return;
				}

				selected = details.items[0] as T;
				onselect?.(details.items[0] as T);
			}
		}),
		{
			context: {
				get value() {
					return [selected[valueKey]];
				}
			}
		}
	);

	const api = $derived(select.connect(snapshot, send, normalizeProps));
</script>

<div {...api.getRootProps()}>
	<select {...api.getHiddenSelectProps()}>
		{#each items as item}
			<option value={item[valueKey]}>
				{item[labelKey]}
			</option>
		{/each}
	</select>

	<div {...api.getControlProps()}>
		<button {...api.getTriggerProps()}>
			<span>
				{api.valueAsString || 'Select Option'}
			</span>

			<span>
				<AltArrowIcon height={18} width={18} />
			</span>
		</button>
	</div>

	<div use:portal {...api.getPositionerProps()}>
		<ul {...api.getContentProps()}>
			{#each items as item (item[valueKey])}
				<li {...api.getItemProps({ item })}>
					<span {...api.getItemTextProps({ item })}>{item[labelKey]}</span>

					<span {...api.getItemIndicatorProps({ item })}>
						<CheckIcon height={18} width={18} />
					</span>
				</li>
			{/each}
		</ul>
	</div>
</div>

<style lang="postcss">
	[data-scope='select'][data-part='root'] {
		[data-part='control'] {
			@apply inline-flex flex-col gap-1;

			[data-part='label'] {
				@apply label-text;
			}

			[data-part='trigger'] {
				@apply border-2 border-neutral-content px-2 py-1 rounded-btn transition-colors flex items-center justify-between gap-3;

				&:focus {
					@apply ring-2 ring-primary/25;
				}

				&:hover {
					@apply border-secondary text-secondary bg-secondary-ghost;
				}

				&:not([data-placeholder-shown]) {
					@apply border-primary text-primary bg-primary-ghost font-medium;
				}

				&[data-placeholder-shown] {
					@apply opacity-85 italic;
				}
			}
		}
	}

	[data-scope='select'][data-part='positioner'] {
		[data-part='content'] {
			@apply min-w-48 bg-base-100 rounded-btn border-2 flex flex-col outline-none border-primary;

			&[data-state='closed'] {
				@apply hidden;
			}

			&:focus {
				@apply ring-2 ring-primary/25;
			}

			[data-part='item'] {
				@apply flex items-center justify-between gap-3 p-2 cursor-pointer transition-colors;

				&[data-highlighted] {
					@apply bg-secondary-ghost text-secondary;
				}

				&[data-state='checked'] {
					@apply bg-primary-ghost text-primary;
				}

				[data-part='item-text'] {
					@apply font-medium text-sm;
				}
			}
		}
	}
</style>
