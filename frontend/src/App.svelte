<script lang="ts">
  import { ProcessFiles, SaveFile } from '../wailsjs/go/main/App';

  let fileNSE: File | null = null;
  let fileBSE: File | null = null;
  let fileCDSL: File | null = null;
  let outputMessage = "";

  async function handleSubmit() {
    if (!fileNSE || !fileBSE || !fileCDSL) {
      outputMessage = "❌ Please upload all 3 files.";
      return;
    }

    const tempDir = "/tmp"; // Alternatively, use: window.api.runtime.TempDir()
    const ts = Date.now();

    const nsePath = `${tempDir}/${ts}_nse.csv`;
    const bsePath = `${tempDir}/${ts}_bse.csv`;
    const cdslPath = `${tempDir}/${ts}_cdsl.csv`;

    try {
      // Save uploaded files to disk
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

      // Process saved files
      const resultPath = await ProcessFiles(nsePath, bsePath, cdslPath);

      // Show download link
      outputMessage = `✅ <a href="file://${resultPath}" download>Download output</a>`;
    } catch (err) {
      outputMessage = `❌ Error: ${err}`;
    }
  }
</script>

<!-- File Upload Inputs -->
<input type="file" on:change={(e) => fileNSE = e.target.files[0]} />
<input type="file" on:change={(e) => fileBSE = e.target.files[0]} />
<input type="file" on:change={(e) => fileCDSL = e.target.files[0]} />

<!-- Submit Button -->
<button on:click={handleSubmit}>Process</button>

<!-- Render Output Message as HTML -->
<p>{@html outputMessage}</p>

