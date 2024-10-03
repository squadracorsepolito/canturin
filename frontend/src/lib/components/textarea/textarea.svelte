<script lang="ts">
	import type { Action } from 'svelte/action';

	type Props = {
		value: string;
		name: string;
		label?: string;
		focusOnDisplay?: boolean;
		rows?: number;
		onescape?: () => void;
	};

	let { value = $bindable(), name, label, focusOnDisplay, rows = 8, onescape }: Props = $props();

	function handleKeydown(e: KeyboardEvent) {
		if (!onescape) return;

		if (e.key === 'Escape') onescape();
	}

	const action: Action<HTMLTextAreaElement> = (el) => {
		if (focusOnDisplay) {
			el.focus();
		}
	};
</script>

<label class="form-control">
	{#if label}
		<div class="label">
			<span class="label-text-alt">
				{label}
			</span>
		</div>
	{/if}

	<textarea
		bind:value
		{name}
		{rows}
		onkeydown={handleKeydown}
		use:action
		class="textarea textarea-primary w-full"
	></textarea>
</label>
