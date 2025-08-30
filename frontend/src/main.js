import './style.css';
import {
  SelectInputFiles,
  DefaultOutputDirectory,
  SelectOutputDirectory,
  RunFinalCD,
  NativeAlert
} from '../wailsjs/go/main/App';

// --- DOM Element Caching ---
const batchFiles = document.getElementById('batchFiles');
const outFile = document.getElementById('outFile');
const outputLog = document.getElementById('outputLog');
const addFilesBtn = document.getElementById('addFiles');
const removeFilesBtn = document.getElementById('removeFiles');
const outDialogBtn = document.getElementById('outDialog');
const convertBtn = document.getElementById('convertButton');
const paramBox = document.getElementById('paramBox');
const gainBox = document.getElementById('gainBox');
const thirtyTwoBitOutput = document.getElementById('thirtyTwoBitOutput');

const controls = [
  batchFiles, outFile, addFilesBtn, removeFilesBtn, outDialogBtn, convertBtn, gainBox, thirtyTwoBitOutput,
  ...document.querySelectorAll('input[name="filter"]'),
  ...document.querySelectorAll('input[name="dither"]')
];

// --- Functions ---

function basename(path) {
  return path.split('/').pop().split('\\').pop();
}

async function setDefaultOutputDir() {
  if (outFile.value !== '') {
    return;
  }
  if (batchFiles.options.length > 0) {
    const firstFile = batchFiles.options[0].value;
    outFile.value = await DefaultOutputDirectory(firstFile);
  }
}

function cliParams() {
  const filter = document.querySelector('input[name="filter"]:checked').value;
  const dither = document.querySelector('input[name="dither"]:checked').value;
  const bit32Arg = thirtyTwoBitOutput.checked ? '/32' : '';
  const gainArg = gainBox.value !== '1.0' ? `/x${gainBox.value}` : '';

  return [bit32Arg, filter, dither, gainArg].filter(item => item !== "" && typeof item === "string");
}

function setControlsEnabled(enabled) {
  controls.forEach(control => control.disabled = !enabled);
}

// --- Event Listeners ---

addFilesBtn.onclick = async () => {
  const files = await SelectInputFiles();
  if (!files || files.length === 0) return;

  batchFiles.innerHTML = '';
  files.forEach(f => {
    const option = document.createElement('option');
    option.value = f;
    option.textContent = f;
    batchFiles.appendChild(option);
  });
  await setDefaultOutputDir();
};

removeFilesBtn.onclick = () => {
  const selectedOptions = Array.from(batchFiles.selectedOptions);
  selectedOptions.forEach(opt => opt.remove());
};

outDialogBtn.onclick = async () => {
  const selectedDirectory = await SelectOutputDirectory();
  if (selectedDirectory && selectedDirectory !== '') {
    outFile.value = selectedDirectory;
  }
};

convertBtn.onclick = async () => {
  const inputFilePaths = Array.from(batchFiles.options).map(option => option.value);
  const outDir = outFile.value;

  if (inputFilePaths.length === 0) {
    await NativeAlert("Bro, you gotta give me something to work with", "Add some input files to process.");
    return;
  }

  if (outDir === '') {
    await NativeAlert("What am I supposed to do with all this?", "Select a directory for the output files.");
    return;
  }

  setControlsEnabled(false);

  try {
    for (const inputFilePath of inputFilePaths) {
      const options = cliParams();
      const inputFileName = basename(inputFilePath);

      outputLog.innerHTML += `<p>Processing <span style="font-weight: bold">${inputFileName}</span>...</p><br>`;
      outputLog.innerHTML += `<p style="font-family: monospace; white-space: pre;">&gt; finalcd.exe <span style="font-style: italic">${cliParams().join(' ')} ${inputFilePath} ${outDir}\\${inputFileName}</span></p><br>`;
      const result = await RunFinalCD(inputFilePath, outDir, options);
      outputLog.innerHTML += `<p style="font-family: monospace;">${result.Stdout.replace(/\r?\n|\r/g, "<br>")}</p>`;
      if (result.Success) {
        outputLog.innerHTML += `<p><span style="font-weight:bold; color:green;">Processed ${result.InputFileName} -> ${result.OutputFilePath}</span></p><hr>`;
      } else {
        outputLog.innerHTML += `<p><span style="font-weight:bold; color:red;">Error processing ${result.InputFileName}: ${result.Err}</span></p><hr>`;
      }
      outputLog.scrollTop = outputLog.scrollHeight;

      const optionToRemove = Array.from(batchFiles.options).find(opt => opt.value === inputFilePath);
      if (optionToRemove) {
        optionToRemove.remove();
      }
    }
  } catch (error) {
    outputLog.innerHTML += `<p><span style="font-weight: bold; color:red;">An unexpected error occurred: ${error.message}</span></p>`;
  } finally {
    outputLog.scrollTop = outputLog.scrollHeight;
    setControlsEnabled(true);
  }
};