import { z } from 'zod';

export function nameSchema(invalidNames: string[]) {
	return z.object({
		name: z
			.string()
			.min(1)
			.refine((n) => !invalidNames.includes(n), { message: 'Duplicated' })
	});
}

export const descSchema = z.object({
	desc: z.string()
});

export const kindSchema = z.number().min(1).max(4);
export const sizeSchema = z.number().min(1).max(64);
export const minSchema = z.number();
export const maxSchema = z.number();
export const scaleSchema = z.number();
export const offsetSchema = z.number();
