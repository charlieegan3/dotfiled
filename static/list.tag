<list>
  <div>
    <div class="file-chunk clearfix" each={ fileChunks }>
		<p class="contents">{ Contents }</p>
		<p class="details">
		  <span class="file-type { FileType }">{ FileType }</span>
		  <span class="count">used <strong>{ Count }</strong> times</span>
		</p>
	</div>
  </div>

  <script>
    this.fileChunks = opts.fileChunks
  </script>
</list>
