<script lang="ts">
	import type { Action } from 'svelte/action';
	import {
		draggable,
		dropTargetForElements
	} from '@atlaskit/pragmatic-drag-and-drop/element/adapter';
	import type { Snippet } from 'svelte';

	type Props = {
		id: string;
		children: Snippet;
	};

	let { id, children }: Props = $props();

	const draggableAction: Action<HTMLDivElement> = (el) => {
		const cleanup = draggable({
			element: el
			// getInitialData() {
			// 	return {
			// 		type: 'item',
			// 		id: id
			// 	};
			// },
		});

		const a = dropTargetForElements({
			element: el,
			onDrop: () => {
				console.log('drop');
			},
			onDragEnter: () => {
				console.log('drag enter');
			},
			onDragLeave: () => {
				console.log('drag leave');
			}
		});

		return {
			destroy() {
				cleanup();
				a();
			}
		};
	};
</script>

<div use:draggableAction>
	{@render children()}
</div>
