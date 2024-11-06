<script lang="ts" generics="T extends {id: string}">
	import type { Snippet } from 'svelte';
	import type { Action } from 'svelte/action';
	import SortableItem from './sortable-item.svelte';
	import { monitorForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import { onMount } from 'svelte';
	import { extractClosestEdge } from '@atlaskit/pragmatic-drag-and-drop-hitbox/closest-edge';
	import { getReorderDestinationIndex } from '@atlaskit/pragmatic-drag-and-drop-hitbox/util/get-reorder-destination-index';
	import { flip } from 'svelte/animate';
	import { isItem, type HighlightState } from './types';

	type Props = {
		items: T[];
		instanceId: string;
		reorder: (id: string, from: number, to: number) => void;
		itemBody: Snippet<[{ item: T; highlightState: HighlightState }]>;
	};

	let { items, instanceId, reorder, itemBody }: Props = $props();

	let selectedItem = $state({
		id: '',
		index: -1
	});

	let mode = $state<'drag' | 'keyboard'>('drag');

	onMount(() => {
		return monitorForElements({
			canMonitor({ source }) {
				return isItem(source.data) && source.data.instanceId === instanceId;
			},
			onDrop({ source, location }) {
				if (location.current.dropTargets.length === 0) return;

				if (!isItem(source.data)) return;

				const itemId = source.data.id;
				// const listId = location.initial.dropTargets[1].data.listId;

				const itemIdx = items.findIndex((item) => item.id === itemId);
				if (itemIdx === -1) return;

				// console.log('onDrop', itemId, itemIdx);

				// if (location.current.dropTargets.length === 1) {
				// 	console.log(
				// 		'dropTargets1',
				// 		location.current.dropTargets,
				// 		location.current.dropTargets.length
				// 	);
				// }

				// if (location.current.dropTargets.length === 1) {
				// Destructure and extract the destination card and column data from the drop targets
				const [destItemRecord] = location.current.dropTargets;
				if (!isItem(destItemRecord.data)) return;

				const destItemId = destItemRecord.data.id;

				// Find the index of the target card within the destination column's cards
				const indexOfTarget = items.findIndex((item) => item.id === destItemId);

				// Determine the closest edge of the target card: top or bottom
				const closestEdgeOfTarget = extractClosestEdge(destItemRecord.data);

				// Calculate the destination index for the card to be reordered within the same column
				const destinationIndex = getReorderDestinationIndex({
					startIndex: itemIdx,
					indexOfTarget,
					closestEdgeOfTarget,
					axis: 'vertical'
				});

				// Perform the card reordering within the same column
				reorder(itemId, itemIdx, destinationIndex);

				return;
				// }
			}
		});
	});

	const listAction: Action<HTMLUListElement> = (el) => {
		function handleKeydown(e: KeyboardEvent) {
			if (e.key === 'Escape') {
				e.preventDefault();
				selectedItem = { id: '', index: -1 };

				mode = 'drag';
				el.blur();
				return;
			}

			if (mode === 'drag') {
				if (e.key === 'Enter') {
					mode = 'keyboard';
					selectedItem.index = 0;
				}
				return;
			}

			if (e.key === 'Enter') {
				selectedItem = { id: '', index: -1 };
				mode = 'drag';

				return;
			}

			if (e.key === ' ') {
				e.preventDefault();

				const targetId = items[selectedItem.index].id;
				if (selectedItem.id === targetId) {
					selectedItem.id = '';
				} else {
					selectedItem.id = targetId;
				}

				return;
			}

			if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
				e.preventDefault();

				if (selectedItem.index <= 0) {
					return;
				}

				if (selectedItem.id) {
					reorder(selectedItem.id, selectedItem.index, selectedItem.index - 1);
				}

				if (selectedItem.index > 0) {
					selectedItem.index--;
				}

				return;
			}

			if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
				e.preventDefault();

				if (selectedItem.index >= items.length - 1) {
					return;
				}

				if (selectedItem.id) {
					reorder(selectedItem.id, selectedItem.index, selectedItem.index + 1);
				}

				if (selectedItem.index < items.length - 1) {
					selectedItem.index++;
				}

				return;
			}
		}

		function handleBlur() {
			selectedItem = { id: '', index: -1 };
			mode = 'drag';
		}

		el.addEventListener('keydown', handleKeydown);
		el.addEventListener('blur', handleBlur);

		return {
			destroy() {
				el.removeEventListener('keydown', handleKeydown);
				el.removeEventListener('blur', handleBlur);
			}
		};
	};
</script>

<!-- svelte-ignore a11y_no_noninteractive_tabindex -->
<ul use:listAction tabindex="0" class="sortable-list">
	{#each items as item, idx (item.id)}
		{@const hState =
			idx === selectedItem.index ? (selectedItem.id === item.id ? 'moving' : 'selecting') : 'none'}

		<li animate:flip={{ duration: 150 }}>
			<SortableItem id={item.id} {instanceId} highlightState={hState}>
				{@render itemBody({ item, highlightState: hState })}
			</SortableItem>
		</li>
	{/each}
</ul>

<style lang="postcss">
	.sortable-list {
		@apply flex flex-col rounded-box;

		&:focus {
			@apply outline-none;
		}
	}
</style>
