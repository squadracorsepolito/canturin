<script lang="ts">
	import { FiniteStateMachine } from 'runed';
	import type { Snippet } from 'svelte';

	type Props = {
		hidePlaceholder?: boolean;
		placeholder: Snippet;
		input: Snippet<[typeof fsm]>;
	};

	let { hidePlaceholder, placeholder, input }: Props = $props();

	type State = 'idle' | 'editing' | 'resetting';
	type Events = 'DBLCLICK' | 'ESCAPE' | 'BLUR' | 'TIMEOUT';
	const fsm = new FiniteStateMachine<State, Events>('idle', {
		idle: {},
		editing: {},
		resetting: {}
	});
</script>

{#if hidePlaceholder && fsm.current === 'idle'}
	<div>
		{@render placeholder()}
	</div>
{/if}

{#if fsm.current === 'editing'}
	<form></form>
{/if}
