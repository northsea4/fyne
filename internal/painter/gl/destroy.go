package gl

// Destroy releases resources owned by this painter. Its GL context must remain
// current until this returns; per-object Free deliberately retains programs.
func (p *painter) Destroy() {
	if p.ctx == nil {
		return
	}
	p.ctx.UseProgram(noProgram)
	release := func(state *programState) {
		if state.uniforms == nil { // not initialized (also supports JS references)
			return
		}
		p.ctx.DeleteBuffer(state.buff)
		p.ctx.DeleteProgram(state.ref)
		*state = programState{}
	}
	for _, state := range []*programState{&p.program, &p.blurProgram, &p.lineProgram,
		&p.rectangleProgram, &p.roundRectangleProgram, &p.polygonProgram,
		&p.arcProgram, &p.bezierCurveProgram, &p.arbitraryPolygonProgram, &p.ellipseProgram} {
		release(state)
	}
	for _, state := range p.shaderPrograms {
		if state.valid {
			release(&state.program)
		}
		for _, texture := range state.textures {
			p.ctx.DeleteTexture(texture.tex)
		}
	}
	p.shaderPrograms = nil
	for _, texture := range p.clippedTextTextures {
		p.ctx.DeleteTexture(texture.texture)
	}
	p.clippedTextTextures = nil
	if p.blurSnapTexValid {
		p.ctx.DeleteTexture(p.blurSnapTex)
		p.blurSnapTexValid = false
	}
	if p.blurKernelTexValid {
		p.ctx.DeleteTexture(p.blurKernelTex)
		p.blurKernelTexValid = false
	}
	// Complete queued use and deletion before the native context disappears.
	p.ctx.Finish()
	p.logError()
	p.ctx, p.canvas, p.contextProvider = nil, nil, nil
}
