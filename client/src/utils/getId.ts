import { Template } from "../types/general";

export function getId(str: string, arr: any[], key: string): number | null {
  for (let i = 0; i < arr.length; i++) {
    if (arr[i][key] === str) {
      return arr[i].id;
    }
  }
  return null;
}

export function getTemplateId(templates: Template[], template: string): number | null {
  templates.forEach((t) => {
    if (t.name === template) {
      return t.id;
    }
  })

  return null;
}
