import { RangoliPattern } from './RangoliPattern';

/**
 * RangoliPattern Usage Examples
 * 
 * Demonstrates various ways to use the RangoliPattern component
 * as a decorative background element.
 */

export const RangoliPatternExamples = () => {
  return (
    <div className="space-y-8 p-8">
      <h1 className="text-3xl font-bold mb-8">RangoliPattern Examples</h1>

      {/* Example 1: Default subtle background */}
      <section>
        <h2 className="text-xl font-semibold mb-4">1. Default Subtle Background</h2>
        <div className="relative bg-white rounded-lg p-8 min-h-[300px] overflow-hidden">
          <RangoliPattern />
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-4">Welcome to VamsaSetu</h3>
            <p className="text-gray-700">
              This card has a subtle rangoli pattern in the background with default opacity (0.1).
              The pattern adds cultural authenticity without distracting from the content.
            </p>
          </div>
        </div>
      </section>

      {/* Example 2: More visible pattern */}
      <section>
        <h2 className="text-xl font-semibold mb-4">2. More Visible Pattern</h2>
        <div className="relative bg-ivory rounded-lg p-8 min-h-[300px] overflow-hidden">
          <RangoliPattern opacity={0.2} />
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-4">Family Tree Dashboard</h3>
            <p className="text-gray-700">
              Increased opacity (0.2) makes the pattern more prominent while still
              maintaining readability of the content.
            </p>
          </div>
        </div>
      </section>

      {/* Example 3: Larger pattern tiles */}
      <section>
        <h2 className="text-xl font-semibold mb-4">3. Larger Pattern Tiles</h2>
        <div className="relative bg-white rounded-lg p-8 min-h-[300px] overflow-hidden">
          <RangoliPattern size={300} opacity={0.15} />
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-4">Events Calendar</h3>
            <p className="text-gray-700">
              Larger pattern tiles (300px) create a more spacious, elegant feel.
              Perfect for hero sections or landing pages.
            </p>
          </div>
        </div>
      </section>

      {/* Example 4: Smaller pattern tiles */}
      <section>
        <h2 className="text-xl font-semibold mb-4">4. Smaller Pattern Tiles</h2>
        <div className="relative bg-ivory rounded-lg p-8 min-h-[300px] overflow-hidden">
          <RangoliPattern size={150} opacity={0.12} />
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-4">Member Profile</h3>
            <p className="text-gray-700">
              Smaller pattern tiles (150px) create a denser, more intricate texture.
              Works well for compact sections or cards.
            </p>
          </div>
        </div>
      </section>

      {/* Example 5: Full page background */}
      <section>
        <h2 className="text-xl font-semibold mb-4">5. Full Page Background</h2>
        <div className="relative bg-white rounded-lg p-8 min-h-[400px] overflow-hidden">
          <RangoliPattern opacity={0.08} size={250} />
          <div className="relative z-10 max-w-2xl mx-auto">
            <h3 className="text-3xl font-bold mb-6 text-center">
              Preserve Your Family Legacy
            </h3>
            <p className="text-gray-700 mb-4">
              VamsaSetu helps you document and celebrate your family's history
              with beautiful, culturally-inspired design elements.
            </p>
            <p className="text-gray-700">
              The rangoli pattern in the background adds a touch of traditional
              Indian artistry while maintaining a modern, clean interface.
            </p>
            <div className="mt-8 flex gap-4 justify-center">
              <button className="px-6 py-3 bg-saffron text-white rounded-lg">
                Get Started
              </button>
              <button className="px-6 py-3 bg-teal text-white rounded-lg">
                Learn More
              </button>
            </div>
          </div>
        </div>
      </section>

      {/* Example 6: With custom className */}
      <section>
        <h2 className="text-xl font-semibold mb-4">6. With Custom Styling</h2>
        <div className="relative bg-gradient-to-br from-ivory to-white rounded-lg p-8 min-h-[300px] overflow-hidden">
          <RangoliPattern 
            opacity={0.15} 
            className="mix-blend-multiply" 
          />
          <div className="relative z-10">
            <h3 className="text-2xl font-bold mb-4">Custom Blend Mode</h3>
            <p className="text-gray-700">
              Using mix-blend-multiply creates interesting interactions with
              gradient backgrounds, adding depth to the design.
            </p>
          </div>
        </div>
      </section>

      {/* Color Reference */}
      <section className="mt-12">
        <h2 className="text-xl font-semibold mb-4">VamsaSetu Theme Colors Used</h2>
        <div className="grid grid-cols-3 gap-4">
          <div className="text-center">
            <div className="w-full h-24 rounded-lg mb-2" style={{ backgroundColor: '#E8650A' }}></div>
            <p className="font-semibold">Saffron</p>
            <p className="text-sm text-gray-600">#E8650A</p>
          </div>
          <div className="text-center">
            <div className="w-full h-24 rounded-lg mb-2" style={{ backgroundColor: '#F5A623' }}></div>
            <p className="font-semibold">Turmeric</p>
            <p className="text-sm text-gray-600">#F5A623</p>
          </div>
          <div className="text-center">
            <div className="w-full h-24 rounded-lg mb-2" style={{ backgroundColor: '#0D4A52' }}></div>
            <p className="font-semibold">Teal</p>
            <p className="text-sm text-gray-600">#0D4A52</p>
          </div>
        </div>
      </section>
    </div>
  );
};

export default RangoliPatternExamples;
