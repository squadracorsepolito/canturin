<script lang="ts">
	import type { BaseEntity, EntityKind } from '$lib/api/canturin';
	import { TextareaEditable, TextEditable } from '$lib/components/editable';
	import { getIconFromEntityKind } from '$lib/components/icon/utils.svelte';
	import { nameSchema, Validator } from '$lib/utils/validator.svelte';
	import { onMount } from 'svelte';

	interface IEntityState {
		entity: BaseEntity;

		getInvalidNames(): Promise<string[]>;
		updateName(name: string): void;
		updateDesc(desc: string): void;
	}

	type Props = {
		entityState: IEntityState;
	};

	let { entityState }: Props = $props();

	const Icon = getIconFromEntityKind(entityState.entity.entityKind);

	let invalidNames = $state<string[]>([]);

	onMount(async () => {
		const res = await entityState.getInvalidNames();
		invalidNames = res;
	});

	const nameValidator = new Validator(
		nameSchema(() => invalidNames),
		() => entityState.entity.name
	);

	function handleName(name: string) {
		entityState.updateName(name);
	}

	function handleDesc(desc: string) {
		console.log(desc);
		entityState.updateDesc(desc);
	}
</script>

{#snippet section(ent: BaseEntity)}
	<div class="flex gap-2 items-center">
		<Icon width="48" height="48" />

		<TextEditable
			bind:value={ent.name}
			name="{ent.entityKind}-name"
			oncommit={handleName}
			errors={nameValidator.errors}
			fontWeight="semibold"
			textSize="lg"
			border="transparent"
		/>
	</div>

	<div class="pt-8">
		<TextareaEditable
			initialValue={ent.desc}
			name="{ent.entityKind}-desc"
			triggerLabel="Add Description"
			onsubmit={handleDesc}
		/>
	</div>
{/snippet}

<section>
	{@render section(entityState.entity)}
</section>
