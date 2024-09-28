<script lang="ts" generics="T">
	import { focusOnDisplay } from '$lib/actions';
	import { FiniteStateMachine } from 'runed';

	type Props = {
		initialAttribute: T;
		onSubmit: (attribute: T) => void;
	};

	let { initialAttribute, onSubmit }: Props = $props();

	let attribute = $state(initialAttribute);

	type State = 'idle' | 'editing' | 'resetting';
	type Events = 'DBLCLICK' | 'BLUR' | 'SUBMIT' | 'ESCAPE' | 'TIMEOUT';
	const fsm = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			DBLCLICK: 'editing'
		},
		editing: {
			SUBMIT: 'idle',
			BLUR: 'idle',
			ESCAPE: 'resetting',
			DBLCLICK: 'editing'
		},
		resetting: {
			_enter: () => {
				fsm.debounce(500, 'TIMEOUT');
			},
			TIMEOUT: 'idle'
		}
	});
</script>

<div ondblclick={() => fsm.send('DBLCLICK')} class="text-2xl">{initialAttribute}</div>

{#if fsm.current === 'editing'}
	<form
		onsubmit={() => {
			onSubmit(attribute);
			fsm.send('SUBMIT');
		}}
	>
		<input
			type="number"
			name="attname"
			bind:value={attribute}
			use:focusOnDisplay
			onblur={() => fsm.send('BLUR')}
			onkeydown={(e) => {
				if (e.key == 'Escape') fsm.send('ESCAPE');
			}}
			class="input input-primary"
			min="1"
			max="64"
		/>
	</form>
{/if}

{fsm.current}
