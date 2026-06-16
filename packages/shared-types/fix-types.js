const fs = require('fs');
const path = require('path');

const typesPath = path.join(__dirname, 'src', 'backend.d.ts');
let content = fs.readFileSync(typesPath, 'utf-8');

// Find all definitions
const regex = /"([a-z_]+\.[A-Za-z0-9_]+)":\s*\{/g;
let match;
let appendStr = '\n// Auto-generated flattened exports for openapi-typescript v5 backwards compatibility\n';

const added = new Set();

while ((match = regex.exec(content)) !== null) {
  const originalKey = match[1];
  // Convert "batches.Batch" -> "BatchesBatch", "categories.Category" -> "CategoriesCategory"
  const parts = originalKey.split('.');
  const prefix = parts[0].charAt(0).toUpperCase() + parts[0].slice(1);
  const suffix = parts[1];
  const flatName = `${prefix}${suffix}`;
  
  // also support some exact matches if prefix is already singular
  // the typical format is domain.Type
  
  if (!added.has(flatName)) {
    appendStr += `export type ${flatName} = components['definitions']['${originalKey}'];\n`;
    added.add(flatName);
  }
}

// Special cases that might be expected by the UI due to naming conventions
appendStr += `
// Aliases just in case
export type Category = components['definitions']['categories.Category'];
export type Batch = components['definitions']['batches.Batch'];
export type Product = components['definitions']['products.Product'];
`;

fs.writeFileSync(typesPath, content + appendStr);
console.log('Appended flattened exports to backend.d.ts');
