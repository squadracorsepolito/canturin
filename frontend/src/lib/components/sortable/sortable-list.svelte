<script lang="ts" generics="T extends {id: string}">
	import { dropTargetForElements } from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import type { Snippet } from 'svelte';
	import type { Action } from 'svelte/action';
	import SortableItem from './sortable-item.svelte';

	type Props = {
		items: T[];
		itemBody: Snippet<[T]>;
	};

	let { items, itemBody }: Props = $props();

	// const listAction: Action<HTMLElement> = (el) => {
	// 	const cleanup = dropTargetForElements({
	// 		element: el,
	// 		getData() {
	// 			return { listId: 'list-1' };
	// 		},
	// 		getIsSticky() {
	// 			return true;
	// 		}
	// 	});

	// 	return {
	// 		destroy() {
	// 			cleanup();
	// 		}
	// 	};
	// };
</script>

<div class="border border-yellow-500 p-5 pb-20 flex flex-col">
	{#each items as item}
		<SortableItem id={item.id}>
			{@render itemBody(item)}
		</SortableItem>
	{/each}
</div>
