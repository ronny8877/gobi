export const guideGroups = [
  {
    label: 'Start here',
    items: [
      ['quick-start', 'Quick start'],
      ['configuration', 'Configuration'],
    ],
  },
  {
    label: 'Build your API',
    items: [
      ['routes', 'Routes & HTTP methods'],
      ['generators', 'Data generators'],
      ['arrays', 'Arrays'],
      ['references', 'Reusable references'],
      ['validation-auth', 'Validation & auth'],
    ],
  },
  {
    label: 'Go further',
    items: [
      ['pocketbase', 'PocketBase persistence'],
      ['deployment', 'Deploy & develop'],
    ],
  },
];
export const guides = guideGroups.flatMap((g) => g.items);
