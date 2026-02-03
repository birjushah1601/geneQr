// Fix corrupted Unicode characters in source files
const fs = require('fs');
const path = require('path');

// Unicode replacements
const replacements = [
  // Corrupted arrows
  ['â†', '←'],
  
  // Corrupted checkmarks  
  ['âœ"', '✓'],
  ['âœ…', '✅'],
  
  // Corrupted emojis
  ['ðŸ'‹', '👋'],
  ['ðŸ"§', '📧'],
  ['ðŸ"±', '📱'],
  ['ðŸš€', '🚀'],
  ['ðŸ'¼', '💼'],
  
  // Corrupted spaces
  ['Â ', ' '],
  ['Â', ''],
  
  // GeneQR rebranding
  ['GeneQR', 'ServQR'],
  ['genq-admin-ui', 'servqr-admin-ui'],
  ['genq', 'servqr'],
];

function walkDir(dir) {
  const files = [];
  const items = fs.readdirSync(dir);
  
  for (const item of items) {
    const fullPath = path.join(dir, item);
    const stat = fs.statSync(fullPath);
    
    if (stat.isDirectory() && item !== 'node_modules' && item !== '.next') {
      files.push(...walkDir(fullPath));
    } else if (stat.isFile() && /\.(tsx?|jsx?|json)$/.test(item)) {
      files.push(fullPath);
    }
  }
  
  return files;
}

const files = walkDir('src');
files.push('package.json');

let totalFixed = 0;
let totalReplacements = 0;

for (const file of files) {
  try {
    let content = fs.readFileSync(file, 'utf8');
    const original = content;
    let fileReplacements = 0;
    
    for (const [from, to] of replacements) {
      const regex = new RegExp(from.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), 'g');
      const matches = content.match(regex);
      if (matches) {
        content = content.replace(regex, to);
        fileReplacements += matches.length;
      }
    }
    
    if (content !== original) {
      fs.writeFileSync(file, content, 'utf8');
      console.log(`✓ Fixed: ${path.relative(process.cwd(), file)} (${fileReplacements} replacements)`);
      totalFixed++;
      totalReplacements += fileReplacements;
    }
  } catch (err) {
    console.error(`Error processing ${file}:`, err.message);
  }
}

console.log('');
console.log('='.repeat(50));
console.log(`Files fixed: ${totalFixed}`);
console.log(`Total replacements: ${totalReplacements}`);
console.log('='.repeat(50));
