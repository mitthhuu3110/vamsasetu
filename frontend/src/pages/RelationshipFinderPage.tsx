import React, { useState } from 'react';
import { useMembers } from '../hooks/useMembers';
import { useFindPath } from '../hooks/useRelationships';
import Card, { CardHeader, CardTitle, CardContent } from '../components/ui/Card';
import Button from '../components/ui/Button';
import LoadingSpinner from '../components/common/LoadingSpinner';

const RelationshipFinderPage: React.FC = () => {
  const [fromMemberId, setFromMemberId] = useState('');
  const [toMemberId, setToMemberId] = useState('');
  const [searchFromId, setSearchFromId] = useState('');
  const [searchToId, setSearchToId] = useState('');

  const { data: membersResponse } = useMembers();
  const members = membersResponse?.data?.members || [];

  const { data: pathResponse, isLoading } = useFindPath(searchFromId, searchToId);

  const handleFindRelationship = () => {
    if (fromMemberId && toMemberId) {
      setSearchFromId(fromMemberId);
      setSearchToId(toMemberId);
    }
  };

  const path = pathResponse?.data;

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Header */}
      <div>
        <h1 className="font-display text-3xl font-bold text-charcoal">
          Relationship Finder
        </h1>
        <p className="text-charcoal/70 mt-1">
          Discover how family members are related to each other
        </p>
      </div>

      {/* Search Form */}
      <Card variant="elevated">
        <CardHeader>
          <CardTitle>Find Relationship</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-charcoal mb-2">
              From Member
            </label>
            <select
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
              value={fromMemberId}
              onChange={(e) => setFromMemberId(e.target.value)}
            >
              <option value="">Select a member</option>
              {members.map((member) => (
                <option key={member.id} value={member.id}>
                  {member.name}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-charcoal mb-2">
              To Member
            </label>
            <select
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-saffron focus:border-transparent"
              value={toMemberId}
              onChange={(e) => setToMemberId(e.target.value)}
            >
              <option value="">Select a member</option>
              {members.map((member) => (
                <option key={member.id} value={member.id}>
                  {member.name}
                </option>
              ))}
            </select>
          </div>

          <Button
            variant="primary"
            onClick={handleFindRelationship}
            fullWidth
            disabled={!fromMemberId || !toMemberId || isLoading}
            isLoading={isLoading}
          >
            {isLoading ? 'Finding...' : 'Find Relationship'}
          </Button>
        </CardContent>
      </Card>

      {/* Results */}
      {searchFromId && searchToId && (
        <Card variant="elevated">
          <CardHeader>
            <CardTitle>Result</CardTitle>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <LoadingSpinner />
            ) : path ? (
              <div className="space-y-4">
                <div className="bg-saffron/10 border border-saffron/20 rounded-lg p-4">
                  <p className="text-charcoal font-medium">
                    {path.description || 'Relationship found!'}
                  </p>
                  {path.relationLabel && (
                    <p className="text-charcoal/70 text-sm mt-1">
                      Relation: {path.relationLabel}
                    </p>
                  )}
                </div>
                {path.path && path.path.length > 0 && (
                  <div>
                    <p className="text-sm font-medium text-charcoal mb-2">Path:</p>
                    <div className="flex flex-wrap items-center gap-2">
                      {path.path.map((step, index) => (
                        <React.Fragment key={index}>
                          <span className="bg-ivory px-3 py-1 rounded-lg text-sm">
                            {typeof step === 'string' ? step : step.name}
                          </span>
                          {index < path.path.length - 1 && (
                            <span className="text-charcoal/40">→</span>
                          )}
                        </React.Fragment>
                      ))}
                    </div>
                  </div>
                )}
              </div>
            ) : (
              <div className="text-center py-8">
                <p className="text-charcoal/60">
                  No relationship path found between these members.
                </p>
              </div>
            )}
          </CardContent>
        </Card>
      )}
    </div>
  );
};

export default RelationshipFinderPage;
