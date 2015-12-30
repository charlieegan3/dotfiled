<list>
  <div>
    <div class="file-chunk clearfix" each={ fileChunks }>
		<p class="contents">{ Contents }</p>
		<p class="details">
		  <span class="file-type">{ FileType }</span>
		  <span class="count">{ Count } files</span>
		</p>
	</div>
  </div>

  <script>
    this.fileChunks = opts.fileChunks
  </script>

  <style scoped>
    :scope {
    }

    .file-chunk {
      margin-bottom: 10px;
    }

    .contents {
      color: #f2f2f2;
      padding: 5px;
      margin: 0px;
      font-family: "Share Tech Mono", monospace;
      font-size: 1.5em;
      text-align: center;
      background-color: #333;
    }

    .details {
      margin-top: 0px;
    }

    .file-type, .count {
      font-size: 0.9em;
      padding: 3px 6px 2px 6px;
      width: 50%;
      box-sizing: border-box;
      border-bottom: 2px solid #333;
    }

    .count {
      border-left: 1px solid #333;
      border-right: 2px solid #333;
      float: right;
    }

    .file-type {
      border-right: 1px solid #333;
      border-left: 2px solid #333;
      float: left;
      text-align: right;
    }
  </style>
</list>
