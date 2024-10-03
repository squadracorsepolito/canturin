<script lang="ts">
	import { Button, IconButton, SubmitButton } from '$lib/components/button';
	import { EditableForm } from '$lib/components/form';
	import { AddIcon, SignalTypeIcon } from '$lib/components/icon';
	import { ResizeableTextInput } from '$lib/components/input';
	import { Textarea } from '$lib/components/textarea';
	import type { PanelSectionProps } from '../types';
	import { getSignalTypeState } from '../../state/signal-type-state.svelte';
	import { descSchema, nameSchema } from './signal-type-schema';
	import { type SignalType } from '$lib/api/canturin';

	let { entityId }: PanelSectionProps = $props();

	let sts = getSignalTypeState(entityId);

	let invalidNames = $state<string[]>([]);

	async function loadInvalidNames() {
		const res = await sts.getInvalidNames();
		invalidNames = res;
	}

	$effect(() => {
		loadInvalidNames();
	});

	function handleName(name: string) {
		sts.updateName(name);
	}

	function handleDesc(desc: string) {
		sts.updateDesc(desc);
	}
</script>

{#snippet section(signalType: SignalType)}
	<div class="flex h-24 gap-3 items-center">
		<SignalTypeIcon width="48" height="48" />

		<div class="flex-1">
			<EditableForm
				schema={nameSchema(invalidNames)}
				initialValues={{ name: signalType.name }}
				hidePlaceholder
				onsubmit={({ name }) => handleName(name)}
			>
				{#snippet placeholder()}
					<div class="text-xl font-bold break-all">{signalType.name}</div>
				{/snippet}

				{#snippet input({ fsm, values, errors })}
					<ResizeableTextInput
						size="sm"
						label="Name"
						name="signal-type-name"
						bind:value={values.name}
						errors={errors.name}
						focusOnDisplay
						onescape={() => fsm.send('ESCAPE')}
						onblur={() => fsm.send('BLUR')}
					/>
				{/snippet}
			</EditableForm>
		</div>
	</div>

	<div>
		<EditableForm
			schema={descSchema}
			initialValues={{ desc: signalType.desc }}
			onsubmit={({ desc }) => handleDesc(desc)}
			hidePlaceholder
			blurOnSubmit
		>
			{#snippet placeholder(fsm)}
				{#if signalType.desc}
					<p>{signalType.desc}</p>
				{:else}
					<IconButton onclick={() => fsm.send('DBLCLICK')} label="Add Description">
						{#snippet icon()}
							<AddIcon />
						{/snippet}
					</IconButton>
				{/if}
			{/snippet}

			{#snippet input({ fsm, values })}
				<Textarea
					bind:value={values.desc}
					name="signal-type-desc"
					onescape={() => fsm.send('ESCAPE')}
					focusOnDisplay
					label="Description"
					rows={20}
				/>

				<div class="flex justify-end gap-3 pt-5">
					<Button label="Cancel" onclick={() => fsm.send('ESCAPE')} />
					<SubmitButton label="Save" />
				</div>
			{/snippet}
		</EditableForm>
	</div>
{/snippet}

<section>
	{@render section(sts.entity)}
</section>
