<script lang="ts">
  import { ProcessFiles, SaveFile, MoveToDownloads } from '../wailsjs/go/main/App';

  let fileNSE: File | null = null;
  let fileBSE: File | null = null;
  let fileCDSL: File | null = null;
  let outputMessage = "";
  let isProcessing = false;

  function handleFileChange(e: Event, setFile: (f: File | null) => void) {
    const target = e.target as HTMLInputElement;
    setFile(target.files?.[0] ?? null);
  }

  async function handleSubmit() {
  if (!fileNSE || !fileBSE || !fileCDSL) {
    outputMessage = "❌ Please upload all 3 files.";
    return;
  }

  isProcessing = true;
  outputMessage = ""; // Clear previous output

  const tempDir = "/tmp";
  const ts = Date.now();

  const nsePath = `${tempDir}/${ts}_nse.csv`;
  const bsePath = `${tempDir}/${ts}_bse.csv`;
  const cdslPath = `${tempDir}/${ts}_cdsl.csv`;

  try {
    await Promise.all([
      fileNSE.arrayBuffer().then(buf =>
        SaveFile(nsePath, Array.from(new Uint8Array(buf)))
      ),
      fileBSE.arrayBuffer().then(buf =>
        SaveFile(bsePath, Array.from(new Uint8Array(buf)))
      ),
      fileCDSL.arrayBuffer().then(buf =>
        SaveFile(cdslPath, Array.from(new Uint8Array(buf)))
      ),
    ]);

    const resultPath = await ProcessFiles(nsePath, bsePath, cdslPath);
    const finalPath = await MoveToDownloads(resultPath);
    outputMessage = `✅ File saved to Downloads folder: ${finalPath}`;
  } catch (err) {
    outputMessage = `❌ Error: ${err}`;
  } finally {
    isProcessing = false;
  }
}

</script>

<style>
  .container {
    max-width: 600px;
    margin: 40px auto;
    padding: 20px;
    font-family: 'Segoe UI', sans-serif;
    color: #fff;
  }

  .file-row {
    display: flex;
    align-items: center;
    margin-bottom: 20px;
  }

  .file-label {
  padding: 10px 16px;
  background-color: #E6E6FA;
  color: #4b0082;
  border-radius: 6px;
  cursor: pointer;
  white-space: nowrap;
  font-size: 1rem;
  flex-shrink: 0;
  transition: background-color 0.2s ease;
  border: 1px solid #c7c7ee;
}

.file-label:hover {
  background-color: #ddd8f5;
}

  .file-input {
    display: none;
  }

  .file-name {
    margin-left: 16px;
    font-style: italic;
    font-size: 0.95rem;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  button {
  padding: 12px 24px;
  font-size: 1rem;
  border: none;
  border-radius: 6px;
  background-color: #A8E6CF;
  color: #034d32; 
  cursor: pointer;
  margin-top: 20px;
  transition: background-color 0.2s ease;
}

button:hover {
  background-color: #94dbc1;
}

  p {
    margin-top: 20px;
    font-size: 1rem;
  }
</style>


<div class="container">

  <!-- NSE File Input -->
  <div class="file-row">
    <label class="file-label">
      Select NSE File (.csv)
      <input class="file-input" type="file" on:change={(e) => handleFileChange(e, f => fileNSE = f)} />
    </label>
    {#if fileNSE}
      <span class="file-name">{fileNSE.name}</span>
    {/if}
  </div>

  <!-- BSE File Input -->
  <div class="file-row">
    <label class="file-label">
      Select BSE File (.csv)
      <input class="file-input" type="file" on:change={(e) => handleFileChange(e, f => fileBSE = f)} />
    </label>
    {#if fileBSE}
      <span class="file-name">{fileBSE.name}</span>
    {/if}
  </div>

  <!-- CDSL File Input -->
  <div class="file-row">
    <label class="file-label">
      Select CDSL File (.csv)
      <input class="file-input" type="file" on:change={(e) => handleFileChange(e, f => fileCDSL = f)} />
    </label>
    {#if fileCDSL}
      <span class="file-name">{fileCDSL.name}</span>
    {/if}
  </div>

  <!-- Submit Button -->
  <button on:click={handleSubmit} disabled={isProcessing || !fileNSE || !fileBSE || !fileCDSL}>
    {isProcessing ? 'Processing...' : 'Process'}
  </button>

  <!-- Output Message -->
  <p>{@html outputMessage}</p>
</div>
