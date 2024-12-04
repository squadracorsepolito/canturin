<script lang="ts">
	import { IconButton } from '$lib/components/button';
	import { UndoIcon, RedoIcon } from '$lib/components/icon';
	import history from '$lib/state/history-state.svelte';
	import { PaneGroup, Pane, PaneResizer } from 'paneforge';
	import { Sidebar } from '$lib/components/sidebar';

	let { children } = $props();

	function handleUndo() {
		history.undo();
	}

	function handleRedo() {
		history.redo();
	}
</script>

<PaneGroup direction="horizontal" class="h-full w-full">
	<Pane defaultSize={15} class="h-full bg-base-200 flex flex-col">
		<div class="h-12 block bg-base-300 sticky top-0"></div>

		<Sidebar />
	</Pane>

	<PaneResizer
		class="h-full w-1 bg-base-200 data-[active=pointer]:bg-accent hover:bg-accent transition-colors delay-75"
	></PaneResizer>

	<Pane class="flex-1 flex flex-col">
		<div class="h-12 bg-base-200 sticky top-0 block">
			<div class="flex items-center h-full px-5 gap-2">
				<IconButton onclick={handleUndo} disabled={!history.canUndo}>
					<UndoIcon />
				</IconButton>

				<IconButton onclick={handleRedo} disabled={!history.canRedo}>
					<RedoIcon />
				</IconButton>
			</div>
		</div>

		{@render children()}
	</Pane>
</PaneGroup>
