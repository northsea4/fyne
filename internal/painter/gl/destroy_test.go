//go:build !js && !android && !ios && !mobile

package gl

import (
	"testing"

	"fyne.io/fyne/v2/canvas"
	"github.com/stretchr/testify/assert"
)

type releaseContext struct {
	context
	programs []Program
	buffers  []Buffer
	textures []Texture
	unbound  bool
	finished int
}

func (c *releaseContext) UseProgram(program Program)    { c.unbound = program == noProgram }
func (c *releaseContext) DeleteProgram(program Program) { c.programs = append(c.programs, program) }
func (c *releaseContext) DeleteBuffer(buffer Buffer)    { c.buffers = append(c.buffers, buffer) }
func (c *releaseContext) DeleteTexture(texture Texture) { c.textures = append(c.textures, texture) }

func (c *releaseContext) Finish()          { c.finished++ }
func (c *releaseContext) GetError() uint32 { return 0 }

func TestPainterDestroyReleasesOwnedResourcesOnce(t *testing.T) {
	ctx := &releaseContext{}
	p := &painter{
		ctx:     ctx,
		program: programState{ref: 11, buff: 21, uniforms: map[string]*uniformState{}},
		shaderPrograms: map[string]*shaderState{"custom": {
			valid: true, program: programState{ref: 12, buff: 22, uniforms: map[string]*uniformState{}},
			textures: map[string]*shaderTexture{"image": {tex: 31}},
		}},
		clippedTextTextures: map[*canvas.Text]clippedTextTexture{{}: {texture: 32}},
		blurSnapTex:         33, blurSnapTexValid: true, blurKernelTex: 34, blurKernelTexValid: true,
	}
	p.Destroy()
	assert.True(t, ctx.unbound)
	assert.Equal(t, 1, ctx.finished)
	assert.ElementsMatch(t, []Program{11, 12}, ctx.programs)
	assert.ElementsMatch(t, []Buffer{21, 22}, ctx.buffers)
	assert.ElementsMatch(t, []Texture{31, 32, 33, 34}, ctx.textures)
	p.Destroy()
	assert.Len(t, ctx.programs, 2)
	assert.Len(t, ctx.buffers, 2)
	assert.Len(t, ctx.textures, 4)
	assert.Equal(t, 1, ctx.finished)
}
