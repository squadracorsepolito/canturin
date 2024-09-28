<script lang="ts">
	import { FiniteStateMachine } from 'runed';
	import { AddIcon } from '../icon';
	import type { Action } from 'svelte/action';

	type Props = {
		initialDesc: string;
		onSubmit: (desc: string) => void;
	};

	let { initialDesc, onSubmit }: Props = $props();

	let desc = $state(initialDesc);

	type State = 'idle' | 'typing' | 'resetting';
	type Events = 'ADD_CLICKED' | 'DBLCLICK' | 'BLUR' | 'ESCAPE' | 'TIMEOUT';

	const f = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			ADD_CLICKED: 'typing',
			DBLCLICK: 'typing'
		},
		typing: {
			BLUR: () => {
				onSubmit(desc);
				return 'idle';
			},
			ESCAPE: 'resetting'
		},
		resetting: {
			_enter: () => {
				f.send('TIMEOUT');
			},
			TIMEOUT: 'idle'
		}
	});

	const textareaAction: Action<HTMLTextAreaElement> = (el) => {
		el.focus();
	};
</script>

{#if f.current === 'idle'}
	{#if initialDesc.length === 0}
		<button onclick={() => f.send('ADD_CLICKED')} class="btn btn-sm btn-ghost">
			<AddIcon />Add Description
		</button>
	{:else}
		<p ondblclick={() => f.send('DBLCLICK')}>{initialDesc}</p>
	{/if}
{:else if f.current === 'typing'}
	<textarea
		bind:value={desc}
		onblur={() => f.current === 'typing' && f.send('BLUR')}
		onkeydown={(e) => {
			if (e.key == 'Escape') f.send('ESCAPE');
		}}
		use:textareaAction
		rows="8"
		class="textarea textarea-primary w-full"
	></textarea>
{/if}
