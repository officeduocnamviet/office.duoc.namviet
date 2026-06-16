const fs = require('fs');
const path = require('path');

function walkDir(dir, callback) {
  fs.readdirSync(dir).forEach(f => {
    let dirPath = path.join(dir, f);
    let isDirectory = fs.statSync(dirPath).isDirectory();
    isDirectory ? walkDir(dirPath, callback) : callback(dirPath);
  });
}

const targetDir = path.join(__dirname, 'src', 'features');

walkDir(targetDir, (filePath) => {
  if (filePath.includes(path.join('api', 'index.ts')) || filePath.endsWith('Api.ts')) {
    let content = fs.readFileSync(filePath, 'utf8');
    let original = content;
    
    // Replace '/api/ with '/
    content = content.replace(/'\/api\//g, "'/");
    
    // Replace `/api/ with `/
    content = content.replace(/`\/api\//g, "`/");
    
    if (content !== original) {
      fs.writeFileSync(filePath, content, 'utf8');
      console.log('Fixed:', filePath);
    }
  }
});
