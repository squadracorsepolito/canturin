<script lang="ts">
	import type { Action } from 'svelte/action';

	type Props = {
		initialValue: string;
		onsubmit: (value: string) => void;
	};

	let { initialValue, onsubmit }: Props = $props();

	let value = $state(initialValue);
	$effect(() => {
		value = initialValue;
	});

	let isEditing = $state(false);
	let canSubmit = $state(true);

	const focus: Action<HTMLElement> = (node) => {
		node.focus();
		return {
			destroy() {
				if (canSubmit && value !== initialValue) {
					onsubmit(value);
				} else {
					value = initialValue;
				}
				canSubmit = true;
			}
		};
	};
</script>

<div>
	{#if isEditing}
		<input
			use:focus
			bind:value
			type="text"
			onkeydown={(e) => {
				if (e.key === 'Enter') {
					isEditing = false;
				} else if (e.key === 'Escape') {
					canSubmit = false;
					isEditing = false;
				}
			}}
			onblur={() => {
				isEditing = false;
			}}
			class="w-full input input-sm"
		/>
	{:else}
		<!-- svelte-ignore a11y_no_static_element_interactions -->
		<span ondblclick={() => (isEditing = true)}>
			{value}
		</span>
	{/if}
</div>
