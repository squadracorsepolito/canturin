import { z } from 'zod';

export function nameSchema(invalidNames: string[]) {
	return z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});
}

export const sizeSchema = z.object({
	size: z.number().min(1).max(64)
});

export const minSchema = z.object({
	min: z.number()
});

export const maxSchema = z.object({
	max: z.number()
});
