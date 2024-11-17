<script lang="ts">
	import { EditorIcon, SettingIcon } from '$lib/components/icon';
	import ToastGroup from '$lib/components/toast/toast-group.svelte';

	import '../app.css';

	let { children } = $props();
</script>

<main class="bg-base-100 flex overflow-y-hidden">
	<aside class="h-full w-12 flex flex-col justify-between bg-base-300">
		<ul>
			<li>
				<a href="/editor" class="btn btn-square btn-ghost rounded-none">
					<EditorIcon height="28" width="28" />
				</a>
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
		{@render children()}
	</div>

	<ToastGroup />
</main>
