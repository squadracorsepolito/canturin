<script lang="ts">
	import { FiniteStateMachine } from 'runed';
	import type { Snippet } from 'svelte';
	import { z, ZodError } from 'zod';

	type Props = {
		label: string;
		prefixName: string;
		initialValue: string;
		invalidNames: string[];
		onSubmit: (name: string) => void;
		children?: Snippet;
	};

	let { label, prefixName, initialValue, invalidNames, onSubmit, children }: Props = $props();

	const schema = z.object({
		name: z
			.string()
			.min(1, { message: `${label}: must contains at least 1 character` })
			.refine((n) => !invalidNames.includes(n), { message: `${label}: is duplicated` })
	});

	const resettingTimeoutMs = 500;

	let width = $state(0);

	let value = $state(initialValue);

	let error = $derived.by(() => {
		try {
			schema.parse({ name: value });
		} catch (err) {
			if (err instanceof ZodError) {
				return err.issues[0].message;
			}
		}
	});

	let canSubmit = $state(true);

	function handleSubmit() {
		if (value !== initialValue && canSubmit) {
			onSubmit(value);
		}
	}

	function handleFormSubmit(e: SubmitEvent) {
		f.send('SUBMIT');

		if (error) {
			e.preventDefault();
			return;
		}

		handleSubmit();
	}

	let inputEl = $state() as HTMLInputElement;

	type State = 'idle' | 'typing' | 'resetting';
	type Events = 'FOCUS' | 'SUBMIT' | 'BLUR' | 'ESCAPE' | 'TIMEOUT';

	const f = new FiniteStateMachine<State, Events>('idle', {
		idle: {
			FOCUS: () => {
				inputEl.select();
				return 'typing';
			}
		},
		typing: {
			SUBMIT: () => {
				if (!error) {
					return 'idle';
				}
			},
			BLUR: () => {
				if (error) {
					return 'resetting';
				}

				handleSubmit();
				return 'idle';
			},
			ESCAPE: () => {
				canSubmit = false;
				inputEl.blur();
				return 'resetting';
			}
		},
		resetting: {
			_enter: () => {
				value = initialValue;
				f.debounce(resettingTimeoutMs, 'TIMEOUT');
			},
			TIMEOUT: () => {
				canSubmit = true;
				return 'idle';
			}
		}
	});

	const handleInputFocus = () => {
		f.send('FOCUS');
	};

	const handleInputBlur = () => {
		f.send('BLUR');
	};

	const handleInputKeydown = (e: KeyboardEvent) => {
		if (e.key === 'Escape') {
			f.send('ESCAPE');
		}
	};

	function getInputClass() {
		switch (f.current) {
			case 'typing':
				if (error) {
					return 'input-error bg-error text-error-content';
				}
				return 'input-primary';
			case 'resetting':
				return 'bg-warning text-warning-content';
			default:
				return '';
		}
	}
</script>

<form class="relative" onsubmit={handleFormSubmit}>
	<label for={prefixName + '_name'} class="hidden">{label}</label>

	<div class="flex items-center gap-3">
		{#if children}
			{@render children()}
		{/if}

		<input
			type="text"
			name={prefixName + '_name'}
			bind:value
			bind:this={inputEl}
			onfocus={handleInputFocus}
			onblur={handleInputBlur}
			onkeydown={handleInputKeydown}
			class="{getInputClass()} input input-sm px-2 font-medium text-xl transition-colors"
			style:width={`${width + 2 + 1}px`}
		/>

		{#if error}
			<span class="text-xs text-error">{error}</span>
		{/if}
	</div>

	<div
		bind:clientWidth={width}
		class="px-2 inline-block absolute left-0 break-all top-12 font-medium text-xl invisible"
	>
		{value}
	</div>
</form>
