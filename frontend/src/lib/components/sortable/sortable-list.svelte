<script lang="ts" generics="T extends {id: string}">
	import type { Snippet } from 'svelte';
	import type { Action } from 'svelte/action';
	import { flip } from 'svelte/animate';
	import { type HighlightState } from './types';
	import { Sortable } from '$lib/actions/sortable.svelte';
	import DragHandleIcon from '../icon/drag-handle-icon.svelte';
	import { Toggle } from '../toggle';
	import SortIcon from '../icon/sort-icon.svelte';
	import './style.css';

	type Props = {
		items: T[];
		instanceId: string;
		reorder: (id: string, from: number, to: number) => void;
		itemBody: Snippet<[{ item: T; highlightState: HighlightState }]>;
	};

	let { items, instanceId, reorder, itemBody }: Props = $props();

	// let selectedItem = $state({
	// 	id: '',
	// 	index: -1
	// });

	// let mode = $state<'drag' | 'keyboard'>('drag');

	const sortable = new Sortable({
		instanceId,
		enabled: true,
		itemsGetter: () => items,
		reorder
	});

	// const listAction: Action<HTMLUListElement> = (el) => {
	// 	function handleKeydown(e: KeyboardEvent) {
	// 		if (e.key === 'Escape') {
	// 			e.preventDefault();
	// 			selectedItem = { id: '', index: -1 };

	// 			mode = 'drag';
	// 			el.blur();
	// 			return;
	// 		}

	// 		if (mode === 'drag') {
	// 			if (e.key === 'Enter') {
	// 				mode = 'keyboard';
	// 				selectedItem.index = 0;
	// 			}
	// 			return;
	// 		}

	// 		if (e.key === 'Enter') {
	// 			selectedItem = { id: '', index: -1 };
	// 			mode = 'drag';

	// 			return;
	// 		}

	// 		if (e.key === ' ') {
	// 			e.preventDefault();

	// 			const targetId = items[selectedItem.index].id;
	// 			if (selectedItem.id === targetId) {
	// 				selectedItem.id = '';
	// 			} else {
	// 				selectedItem.id = targetId;
	// 			}

	// 			return;
	// 		}

	// 		if (e.key === 'ArrowUp' || e.key === 'ArrowLeft') {
	// 			e.preventDefault();

	// 			if (selectedItem.index <= 0) {
	// 				return;
	// 			}

	// 			if (selectedItem.id) {
	// 				reorder(selectedItem.id, selectedItem.index, selectedItem.index - 1);
	// 			}

	// 			if (selectedItem.index > 0) {
	// 				selectedItem.index--;
	// 			}

	// 			return;
	// 		}

	// 		if (e.key === 'ArrowDown' || e.key === 'ArrowRight') {
	// 			e.preventDefault();

	// 			if (selectedItem.index >= items.length - 1) {
	// 				return;
	// 			}

	// 			if (selectedItem.id) {
	// 				reorder(selectedItem.id, selectedItem.index, selectedItem.index + 1);
	// 			}

	// 			if (selectedItem.index < items.length - 1) {
	// 				selectedItem.index++;
	// 			}

	// 			return;
	// 		}
	// 	}

	// 	function handleBlur() {
	// 		selectedItem = { id: '', index: -1 };
	// 		mode = 'drag';
	// 	}

	// 	el.addEventListener('keydown', handleKeydown);
	// 	el.addEventListener('blur', handleBlur);

	// 	return {
	// 		destroy() {
	// 			el.removeEventListener('keydown', handleKeydown);
	// 			el.removeEventListener('blur', handleBlur);
	// 		}
	// 	};
	// };
</script>

<!-- <ul use:sortableWrapper.root class="sortable-list">
	{#each items as item, idx (item.id)}
		{@const hState =
			idx === sortableWrapper.selectedItem.index
				? sortableWrapper.selectedItem.id === item.id
					? 'moving'
					: 'selecting'
				: 'none'}

		<li animate:flip={{ duration: 150 }}>
			<SortableItem id={item.id} {instanceId} highlightState={hState}>
				{@render itemBody({ item, highlightState: hState })}
			</SortableItem>
		</li>
	{/each}
</ul> -->

<Toggle bind:toggled={sortable.enabled} name="sortable-list-enable">
	<SortIcon />
</Toggle>

<div use:sortable.root class="sortable">
	<div use:sortable.dropIndicator></div>
	{#each items as item (item.id)}
		<div
			animate:flip={{ duration: 150 }}
			use:sortable.item={{
				id: item.id
			}}
		>
			<span use:sortable.dragHandle={{ id: item.id }}><DragHandleIcon /></span>

			{@render itemBody({ item, highlightState: 'none' })}
		</div>
	{/each}
</div>

<!-- <style lang="postcss">
	.sortable-list {
		@apply flex flex-col rounded-box;

		&:focus {
			@apply outline-none;
		}
	}
</style> -->
