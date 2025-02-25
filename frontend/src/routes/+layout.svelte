<script lang="ts">
	import '../app.css';
	import { EditorIcon, SettingIcon, UndoIcon, RedoIcon } from '$lib/components/icon';
	import { ToastGroup } from '$lib/components/toast';
	import { IconButton } from '$lib/components/button';
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

<main class="bg-base-100 flex overflow-y-hidden">
	<aside class="h-full w-12 flex flex-col justify-between bg-base-300">
		<ul>
			<li>
				<button class="btn btn-square btn-ghost rounded-none">
					<EditorIcon height="28" width="28" />
				</button>
			</li>
		</ul>
		<ul class="flex flex-col items-center gap-4">
			<li>
				<button class="btn btn-square btn-ghost rounded-none">
					<SettingIcon height="28" width="28" />
				</button>
			</li>
			<input type="checkbox" value="dark" class="toggle theme-controller toggle-sm" />
		</ul>
	</aside>

	<div class="flex-1">
		<PaneGroup direction="horizontal" class="h-full w-full">
			<Pane defaultSize={15} minSize={5} maxSize={25} class="h-full bg-base-200 flex flex-col">
				<div class="h-12 block bg-base-300 sticky top-0"></div>

				<Sidebar />
			</Pane>

			<PaneResizer
				class="h-full w-1 bg-base-300 data-[active=pointer]:bg-accent hover:bg-accent transition-colors delay-75"
			/>

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
	</div>

	<ToastGroup />
</main>
