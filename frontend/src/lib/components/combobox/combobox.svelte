<script lang="ts" generics="T extends { [K in keyof T]: any }, V extends string">
	import { uniqueId, type FieldNameOf } from '$lib/utils';
	import * as combobox from '@zag-js/combobox';
	import { normalizeProps, useMachine } from '@zag-js/svelte';
	import { AltArrowIcon, CheckIcon, CloseIcon } from '../icon';

	type Props = {
		items: T[];
		selected: V;
		name: string;
		valueKey: FieldNameOf<T, V>;
		labelKey: FieldNameOf<T, string>;
		onselect?: (value: V) => void;
		onitemselect?: (item: T) => void;
	};

	let {
		items,
		selected = $bindable(),
		name,
		valueKey,
		labelKey,
		onselect,
		onitemselect
	}: Props = $props();

	let options = $state.raw(items);

	const collection = combobox.collection({
		items,
		itemToString: (item) => item[labelKey],
		itemToValue: (item) => item[valueKey]
	});

	const [snapshot, send] = useMachine(
		combobox.machine({
			id: uniqueId(),
			name,
			collection,
			inputBehavior: 'autocomplete',
			selectionBehavior: 'replace',
			onOpenChange: () => {
				options = items;
			},
			onInputValueChange: (detalis) => {
				const value = detalis.inputValue;
				const filtered = items.filter((item) =>
					item[labelKey].toLowerCase().includes(value.toLowerCase())
				);

				if (filtered.length === 0) {
					collection.setItems([]);
					options = [];
					return;
				}

				collection.setItems(filtered);
				options = filtered;
			},
			onValueChange: (details) => {
				if (details.items.length === 0) {
					return;
				}

				console.log(details);

				const item = details.items[0] as T;
				const value = details.value[0] as V;

				selected = value;
				onselect?.(value);
				onitemselect?.(item);
			}
		}),
		{
			context: {
				get value() {
					return [selected];
				}
			}
		}
	);

	const api = $derived(combobox.connect(snapshot, send, normalizeProps));
</script>

<div {...api.getRootProps()}>
	<div {...api.getControlProps()}>
		<div
			data-scope="combobox"
			data-part="control-group"
			data-selected={api.hasSelectedItems ? true : undefined}
		>
			<input {...api.getInputProps()} />

			<button {...api.getTriggerProps()}>
				<AltArrowIcon height={18} width={18} />
			</button>
		</div>

		<button {...api.getClearTriggerProps()}>
			<CloseIcon />
		</button>
	</div>

	<div {...api.getPositionerProps()}>
		{#if options.length > 0}
			<ul {...api.getContentProps()}>
				{#each options as item}
					<li {...api.getItemProps({ item })}>
						<span {...api.getItemTextProps({ item })}>{item[labelKey]}</span>

						<span {...api.getItemIndicatorProps({ item })}>
							<CheckIcon height={18} width={18} />
						</span>
					</li>
				{/each}
			</ul>
		{/if}
	</div>
</div>

<style lang="postcss">
	[data-scope='combobox'][data-part='root'] {
		@apply inline-block;

		[data-part='control'] {
			@apply flex items-center relative;

			[data-part='control-group'] {
				@apply border-2 border-neutral-content rounded-btn px-2 py-1 flex items-center gap-3;

				[data-part='input'] {
					@apply outline-none font-medium bg-transparent;
				}

				&[data-selected] {
					@apply border-primary bg-primary-ghost text-primary;
				}
			}

			&[data-focus] {
				[data-part='control-group'] {
					@apply focus-ring-primary border-primary bg-primary-ghost text-primary;
				}
			}

			[data-part='clear-trigger'] {
				@apply text-error absolute -right-8;
			}
		}

		[data-part='positioner'] {
			[data-part='content'] {
				@apply bg-base-100 rounded-btn border-2 border-primary;

				[data-part='item'] {
					@apply p-2 cursor-pointer transition-colors flex items-center gap-3 justify-between;

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
	}
</style>
